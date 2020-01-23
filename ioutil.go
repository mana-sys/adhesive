package adhesive

//
//import (
//	"errors"
//	"os"
//)
//
//// mustDir ensures that the specified directory exists. The directory is
//// created if it does not exist. If the specified path exists but is not
//// a directory, returns ErrNotDir.
//func mustDir(dir string) error {
//	stat, err := os.Stat(dir)
//	if os.IsNotExist(err) {
//		return os.Mkdir(dir, 0755)
//	}
//	if err != nil {
//		return err
//	}
//
//	if !stat.IsDir() {
//		return errors.New("not a directory")
//	}
//
//	return nil
//}
