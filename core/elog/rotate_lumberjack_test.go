package elog

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// !!!NOTE!!!
//
// Running these tests in parallel will almost certainly cause sporadic (or even
// regular) failures, because they're all messing with the same global variable
// that controls the logic's mocked time.Now.  So... don't do that.

// Since all the tests uses the time to determine filenames etc, we need to
// control the wall clock as much as possible, which means having a wall clock
// that doesn't change unless we want it to.
var fakeCurrentTime = time.Now()

func fakeTime() time.Time {
	return fakeCurrentTime
}

//func TestNewFile(t *testing.T) {
//	currentTime = fakeTime
//
//	dir := makeTempDir("TestNewFile", t)
//	defer os.RemoveAll(dir)
//	l := &rLogger{
//		Filename: logFile(dir),
//	}
//	defer l.Close()
//	b := []byte("boo!")
//	n, err := l.Write(b)
//	assert.Nil(t, err)
//	assert.Equal(t, len(b), n)
//	existsWithContent(logFile(dir), b, t)
//	fileCount(dir, 1, t)
//}

func TestOpenExisting(t *testing.T) {
	currentTime = fakeTime
	dir := makeTempDir("TestOpenExisting", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	data := []byte("foo!")
	err := ioutil.WriteFile(filename, data, 0644)
	assert.Nil(t, err)

	existsWithContent(filename, data, t)

	l := &rLogger{
		Filename: filename,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)

	// make sure the file got appended
	existsWithContent(filename, append(data, b...), t)

	// make sure no other files were created
	fileCount(dir, 1, t)
}

func TestWriteTooLong(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1
	dir := makeTempDir("TestWriteTooLong", t)
	defer os.RemoveAll(dir)
	l := &rLogger{
		Filename: logFile(dir),
		MaxSize:  5,
	}
	defer l.Close()
	b := []byte("booooooooooooooo!")
	n, err := l.Write(b)
	assert.NotNil(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, fmt.Sprintf("write length %d exceeds maximum file size %d", len(b), l.MaxSize), err.Error())
	_, err = os.Stat(logFile(dir))
	assert.Equal(t, true, os.IsNotExist(err))
}

func TestMakeLogDir(t *testing.T) {
	currentTime = fakeTime
	dir := time.Now().Format("TestMakeLogDir" + backupTimeFormat)
	dir = filepath.Join(os.TempDir(), dir)
	defer os.RemoveAll(dir)
	filename := logFile(dir)
	l := &rLogger{
		Filename: filename,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)
	existsWithContent(logFile(dir), b, t)
	fileCount(dir, 1, t)
}

func TestDefaultFilename(t *testing.T) {
	currentTime = fakeTime
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	filename := filepath.Join(dir, filepath.Base("lumberjack.log"))
	defer os.Remove(filename)
	l := &rLogger{
		Filename: filename,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)
	existsWithContent(filename, b, t)
}

func TestAutoRotate(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestAutoRotate", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Filename: filename,
		MaxSize:  10,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)

	existsWithContent(filename, b, t)
	fileCount(dir, 1, t)

	newFakeTime()

	b2 := []byte("foooooo!")
	n, err = l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, len(b2), n)

	// the old logfile should be moved aside and the main logfile should have
	// only the last write in it.
	existsWithContent(filename, b2, t)

	// the backup file will use the current fake time and have the old contents.
	existsWithContent(backupFile(dir), b, t)

	fileCount(dir, 2, t)
}

func TestFirstWriteRotate(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1
	dir := makeTempDir("TestFirstWriteRotate", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Filename: filename,
		MaxSize:  10,
	}
	defer l.Close()

	start := []byte("boooooo!")
	err := ioutil.WriteFile(filename, start, 0600)
	assert.Nil(t, err)

	newFakeTime()

	// this would make us rotate
	b := []byte("fooo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)

	existsWithContent(filename, b, t)
	existsWithContent(backupFile(dir), start, t)

	fileCount(dir, 2, t)
}

func TestMaxBackups(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1
	dir := makeTempDir("TestMaxBackups", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 1,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)

	existsWithContent(filename, b, t)
	fileCount(dir, 1, t)

	newFakeTime()

	// this will put us over the max
	b2 := []byte("foooooo!")
	n, err = l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, len(b2), n)

	// this will use the new fake time
	secondFilename := backupFile(dir)
	existsWithContent(secondFilename, b, t)

	// make sure the old file still exists with the same content.
	existsWithContent(filename, b2, t)

	fileCount(dir, 2, t)

	newFakeTime()

	// this will make us rotate again
	b3 := []byte("baaaaaar!")
	n, err = l.Write(b3)
	assert.Nil(t, err)
	assert.Equal(t, len(b3), n)

	// this will use the new fake time
	thirdFilename := backupFile(dir)
	existsWithContent(thirdFilename, b2, t)

	existsWithContent(filename, b3, t)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(time.Millisecond * 10)

	// should only have two files in the dir still
	fileCount(dir, 2, t)

	// second file name should still exist
	existsWithContent(thirdFilename, b2, t)

	// should have deleted the first backup
	notExist(secondFilename, t)

	// now test that we don't delete directories or non-logfile files

	newFakeTime()

	// create a file that is close to but different from the logfile name.
	// It shouldn't get caught by our deletion filters.
	notlogfile := logFile(dir) + ".foo"
	err = ioutil.WriteFile(notlogfile, []byte("data"), 0644)
	assert.Nil(t, err)

	// Make a directory that exactly matches our log file filters... it still
	// shouldn't get caught by the deletion filter since it's a directory.
	notlogfiledir := backupFile(dir)
	err = os.Mkdir(notlogfiledir, 0700)
	assert.Nil(t, err)

	newFakeTime()

	// this will use the new fake time
	fourthFilename := backupFile(dir)

	// Create a log file that is/was being compressed - this should
	// not be counted since both the compressed and the uncompressed
	// log files still exist.
	compLogFile := fourthFilename + compressSuffix
	err = ioutil.WriteFile(compLogFile, []byte("compress"), 0644)
	assert.Nil(t, err)

	// this will make us rotate again
	b4 := []byte("baaaaaaz!")
	n, err = l.Write(b4)
	assert.Nil(t, err)
	assert.Equal(t, len(b4), n)

	existsWithContent(fourthFilename, b3, t)
	existsWithContent(fourthFilename+compressSuffix, []byte("compress"), t)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(time.Millisecond * 10)

	// We should have four things in the directory now - the 2 log files, the
	// not log file, and the directory
	fileCount(dir, 5, t)

	// third file name should still exist
	existsWithContent(filename, b4, t)

	existsWithContent(fourthFilename, b3, t)

	// should have deleted the first filename
	notExist(thirdFilename, t)

	// the not-a-logfile should still exist
	exists(notlogfile, t)

	// the directory
	exists(notlogfiledir, t)
}

func TestCleanupExistingBackups(t *testing.T) {
	// test that if we start with more backup files than we're supposed to have
	// in total, that extra ones get cleaned up when we rotate.

	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestCleanupExistingBackups", t)
	defer os.RemoveAll(dir)

	// make 3 backup files

	data := []byte("data")
	backup := backupFile(dir)
	err := ioutil.WriteFile(backup, data, 0644)
	assert.Nil(t, err)

	newFakeTime()

	backup = backupFile(dir)
	err = ioutil.WriteFile(backup+compressSuffix, data, 0644)
	assert.Nil(t, err)

	newFakeTime()

	backup = backupFile(dir)
	err = ioutil.WriteFile(backup, data, 0644)
	assert.Nil(t, err)

	// now create a primary log file with some data
	filename := logFile(dir)
	err = ioutil.WriteFile(filename, data, 0644)
	assert.Nil(t, err)

	l := &rLogger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 1,
	}
	defer l.Close()

	newFakeTime()

	b2 := []byte("foooooo!")
	n, err := l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, len(b2), n)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(time.Millisecond * 10)

	// now we should only have 2 files left - the primary and one backup
	fileCount(dir, 2, t)
}

func TestMaxAge(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestMaxAge", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Filename: filename,
		MaxSize:  10,
		MaxAge:   1,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, len(b), n)

	existsWithContent(filename, b, t)
	fileCount(dir, 1, t)

	// two days later
	newFakeTime()

	b2 := []byte("foooooo!")
	n, err = l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, len(b2), n)
	existsWithContent(backupFile(dir), b, t)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(10 * time.Millisecond)

	// We should still have 2 log files, since the most recent backup was just
	// created.
	fileCount(dir, 2, t)

	existsWithContent(filename, b2, t)

	// we should have deleted the old file due to being too old
	existsWithContent(backupFile(dir), b, t)

	// two days later
	newFakeTime()

	b3 := []byte("baaaaar!")
	n, err = l.Write(b3)
	assert.Nil(t, err)
	assert.Equal(t, len(b3), n)
	existsWithContent(backupFile(dir), b2, t)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(10 * time.Millisecond)

	// We should have 2 log files - the main log file, and the most recent
	// backup.  The earlier backup is past the cutoff and should be gone.
	fileCount(dir, 2, t)

	existsWithContent(filename, b3, t)

	// we should have deleted the old file due to being too old
	existsWithContent(backupFile(dir), b2, t)
}

func TestOldLogFiles(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestOldLogFiles", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	data := []byte("data")
	err := ioutil.WriteFile(filename, data, 07)
	assert.Nil(t, err)

	// This gives us a time with the same precision as the time we get from the
	// timestamp in the name.
	t1, err := time.Parse(backupTimeFormat, fakeTime().UTC().Format(backupTimeFormat))
	assert.Nil(t, err)

	backup := backupFile(dir)
	err = ioutil.WriteFile(backup, data, 07)
	assert.Nil(t, err)

	newFakeTime()

	t2, err := time.Parse(backupTimeFormat, fakeTime().UTC().Format(backupTimeFormat))
	assert.Nil(t, err)

	backup2 := backupFile(dir)
	err = ioutil.WriteFile(backup2, data, 07)
	assert.Nil(t, err)

	l := &rLogger{Filename: filename}
	files, err := l.oldLogFiles()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(files))

	// should be sorted by newest file first, which would be t2
	assert.Equal(t, t1, files[0].timestamp)
	assert.Equal(t, t2, files[1].timestamp)

}

func TestTimeFromName(t *testing.T) {
	l := &rLogger{Filename: "/var/log/myfoo/foo.log"}
	prefix, ext := l.prefixAndExt()

	tests := []struct {
		filename string
		want     time.Time
		wantErr  bool
	}{
		{"foo-2014-05-04T14-44-33.555.log", time.Date(2014, 5, 4, 14, 44, 33, 555000000, time.UTC), false},
		{"foo-2014-05-04T14-44-33.555", time.Time{}, true},
		{"2014-05-04T14-44-33.555.log", time.Time{}, true},
		{"foo.log", time.Time{}, true},
	}

	for _, test := range tests {
		got, err := l.timeFromName(test.filename, prefix, ext)
		assert.Equal(t, test.want, got)
		assert.Equal(t, test.wantErr, err)
	}
}

func TestLocalTime(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestLocalTime", t)
	defer os.RemoveAll(dir)

	l := &rLogger{
		Filename:  logFile(dir),
		MaxSize:   10,
		LocalTime: true,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b))

	b2 := []byte("fooooooo!")
	n2, err := l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, n2, len(b2))

	existsWithContent(logFile(dir), b2, t)
	existsWithContent(backupFileLocal(dir), b, t)
}

func TestRotate(t *testing.T) {
	currentTime = fakeTime
	dir := makeTempDir("TestRotate", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)

	l := &rLogger{
		Filename:   filename,
		MaxBackups: 1,
		MaxSize:    100, // megabytes
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b))

	existsWithContent(filename, b, t)
	fileCount(dir, 1, t)

	newFakeTime()

	err = l.Rotate()
	assert.Nil(t, err)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(10 * time.Millisecond)

	filename2 := backupFile(dir)
	existsWithContent(filename2, b, t)
	existsWithContent(filename, []byte{}, t)
	fileCount(dir, 2, t)
	newFakeTime()

	err = l.Rotate()
	assert.Nil(t, err)

	// we need to wait a little bit since the files get deleted on a different
	// goroutine.
	<-time.After(10 * time.Millisecond)

	filename3 := backupFile(dir)
	existsWithContent(filename3, []byte{}, t)
	existsWithContent(filename, []byte{}, t)
	fileCount(dir, 2, t)

	b2 := []byte("foooooo!")
	n, err = l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b2))

	// this will use the new fake time
	existsWithContent(filename, b2, t)
}

func TestCompressOnRotate(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestCompressOnRotate", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Compress: true,
		Filename: filename,
		MaxSize:  10,
	}
	defer l.Close()
	b := []byte("boo!")
	n, err := l.Write(b)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b))

	existsWithContent(filename, b, t)
	fileCount(dir, 1, t)

	newFakeTime()

	err = l.Rotate()
	assert.Nil(t, err)

	// the old logfile should be moved aside and the main logfile should have
	// nothing in it.
	existsWithContent(filename, []byte{}, t)

	// we need to wait a little bit since the files get compressed on a different
	// goroutine.
	<-time.After(300 * time.Millisecond)

	// a compressed version of the log file should now exist and the original
	// should have been removed.
	bc := new(bytes.Buffer)
	gz := gzip.NewWriter(bc)
	_, err = gz.Write(b)
	assert.Nil(t, err)

	err = gz.Close()
	assert.Nil(t, err)

	existsWithContent(backupFile(dir)+compressSuffix, bc.Bytes(), t)
	notExist(backupFile(dir), t)

	fileCount(dir, 2, t)
}

func TestCompressOnResume(t *testing.T) {
	currentTime = fakeTime
	megabyte = 1

	dir := makeTempDir("TestCompressOnResume", t)
	defer os.RemoveAll(dir)

	filename := logFile(dir)
	l := &rLogger{
		Compress: true,
		Filename: filename,
		MaxSize:  10,
	}
	defer l.Close()

	// Create a backup file and empty "compressed" file.
	filename2 := backupFile(dir)
	b := []byte("foo!")
	err := ioutil.WriteFile(filename2, b, 0644)
	assert.Nil(t, err)

	err = ioutil.WriteFile(filename2+compressSuffix, []byte{}, 0644)
	assert.Nil(t, err)

	newFakeTime()

	b2 := []byte("boo!")
	n, err := l.Write(b2)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b2))
	existsWithContent(filename, b2, t)

	// we need to wait a little bit since the files get compressed on a different
	// goroutine.
	<-time.After(300 * time.Millisecond)

	// The write should have started the compression - a compressed version of
	// the log file should now exist and the original should have been removed.
	bc := new(bytes.Buffer)
	gz := gzip.NewWriter(bc)
	_, err = gz.Write(b)
	assert.Nil(t, err)

	err = gz.Close()
	assert.Nil(t, err)

	existsWithContent(filename2+compressSuffix, bc.Bytes(), t)
	notExist(filename2, t)

	fileCount(dir, 2, t)
}

func TestJson(t *testing.T) {
	data := []byte(`
{
	"filename": "foo",
	"maxsize": 5,
	"maxage": 10,
	"maxbackups": 3,
	"localtime": true,
	"compress": true
}`[1:])

	l := rLogger{}
	err := json.Unmarshal(data, &l)
	assert.Nil(t, err)
	assert.Equal(t, "foo", l.Filename)
	assert.Equal(t, 5, l.MaxSize)
	assert.Equal(t, 10, l.MaxAge)
	assert.Equal(t, 3, l.MaxBackups)
	assert.Equal(t, true, l.LocalTime)
	assert.Equal(t, true, l.Compress)
}

func TestYaml(t *testing.T) {
	data := []byte(`
filename: foo
maxsize: 5
maxage: 10
maxbackups: 3
localtime: true
compress: true`[1:])

	l := rLogger{}
	err := yaml.Unmarshal(data, &l)
	assert.Nil(t, err)
	assert.Equal(t, "foo", l.Filename)
	assert.Equal(t, 5, l.MaxSize)
	assert.Equal(t, 10, l.MaxAge)
	assert.Equal(t, 3, l.MaxBackups)
	assert.Equal(t, true, l.LocalTime)
	assert.Equal(t, true, l.Compress)
}

func TestToml(t *testing.T) {
	data := `
filename = "foo"
maxsize = 5
maxage = 10
maxbackups = 3
localtime = true
compress = true`[1:]

	l := rLogger{}
	md, err := toml.Decode(data, &l)
	assert.Nil(t, err)
	assert.Equal(t, "foo", l.Filename)
	assert.Equal(t, 5, l.MaxSize)
	assert.Equal(t, 10, l.MaxAge)
	assert.Equal(t, 3, l.MaxBackups)
	assert.Equal(t, true, l.LocalTime)
	assert.Equal(t, true, l.Compress)
	assert.Equal(t, 0, len(md.Undecoded()))
}

// makeTempDir creates a file with a semi-unique name in the OS temp directory.
// It should be based on the name of the test, to keep parallel tests from
// colliding, and must be cleaned up after the test is finished.
func makeTempDir(name string, t testing.TB) string {
	//dir := time.Now().Format(name + backupTimeFormat)
	//dir = filepath.Join(os.TempDir(), dir)
	dir, err := ioutil.TempDir(name, "")
	require.Nil(t, err)
	isNilUp(os.Mkdir(dir, 0700), t, 1)
	return dir
}

// existsWithContent checks that the given file exists and has the correct content.
func existsWithContent(path string, content []byte, t testing.TB) {
	info, err := os.Stat(path)
	isNilUp(err, t, 1)
	equalsUp(int64(len(content)), info.Size(), t, 1)

	b, err := ioutil.ReadFile(path)
	isNilUp(err, t, 1)
	equalsUp(content, b, t, 1)
}

// logFile returns the log file name in the given directory for the current fake
// time.
func logFile(dir string) string {
	return filepath.Join(dir, "foobar.log")
}

func backupFile(dir string) string {
	return filepath.Join(dir, "foobar-"+fakeTime().UTC().Format(backupTimeFormat)+".log")
}

func backupFileLocal(dir string) string {
	return filepath.Join(dir, "foobar-"+fakeTime().Format(backupTimeFormat)+".log")
}

// logFileLocal returns the log file name in the given directory for the current
// fake time using the local timezone.
func logFileLocal(dir string) string {
	return filepath.Join(dir, fakeTime().Format(backupTimeFormat))
}

// fileCount checks that the number of files in the directory is exp.
func fileCount(dir string, exp int, t testing.TB) {
	files, err := ioutil.ReadDir(dir)
	isNilUp(err, t, 1)
	// Make sure no other files were created.
	equalsUp(exp, len(files), t, 1)
}

// newFakeTime sets the fake "current time" to two days later.
func newFakeTime() {
	fakeCurrentTime = fakeCurrentTime.Add(time.Hour * 24 * 2)
}

func notExist(path string, t testing.TB) {
	_, err := os.Stat(path)
	assertUp(os.IsNotExist(err), t, 1, "expected to get os.IsNotExist, but instead got %v", err)
}

func exists(path string, t testing.TB) {
	_, err := os.Stat(path)
	assertUp(err == nil, t, 1, "expected file to exist, but got error from os.Stat: %v", err)
}

// isNilUp is like isNil, but used inside helper functions, to ensure that the
// file and line number reported by failures corresponds to one or more levels
// up the stack.
func isNilUp(obtained interface{}, t testing.TB, caller int) {
	if !_isNil(obtained) {
		_, file, line, _ := runtime.Caller(caller + 1)
		fmt.Printf("%s:%d: expected nil, got: %v\n", filepath.Base(file), line, obtained)
		t.FailNow()
	}
}

// _isNil is a helper function for isNil and notNil, and should not be used
// directly.
func _isNil(obtained interface{}) bool {
	if obtained == nil {
		return true
	}

	switch v := reflect.ValueOf(obtained); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

// equalsUp is like equals, but used inside helper functions, to ensure that the
// file and line number reported by failures corresponds to one or more levels
// up the stack.
func equalsUp(exp, act interface{}, t testing.TB, caller int) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(caller + 1)
		fmt.Printf("%s:%d: exp: %v (%T), got: %v (%T)\n",
			filepath.Base(file), line, exp, exp, act, act)
		t.FailNow()
	}
}

// assertUp is like assert, but used inside helper functions, to ensure that
// the file and line number reported by failures corresponds to one or more
// levels up the stack.
func assertUp(condition bool, t testing.TB, caller int, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(caller + 1)
		v = append([]interface{}{filepath.Base(file), line}, v...)
		fmt.Printf("%s:%d: "+msg+"\n", v...)
		t.FailNow()
	}
}
