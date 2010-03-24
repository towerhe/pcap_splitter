package main

import (
  "log"
	"flag"
  "fmt"
	"os"
)

var toBeSplitted = flag.String("f", "", "file to be splitted")
var sizeOfSeg = flag.Int("s", 10 * 1024, "size of segments that file will be splitted into")
//var numOfSeg = flag.Int("n", 10, "number of segments that file will be splitted into")
var dstDir = flag.String("d", "/tmp", "directory path used to save the splitted segments. If the directory path is not existed, it will be created.")

var globalHeader [24]byte

func loadGlobalHeader(file *os.File) (n int, err os.Error) {
    return file.Read(&globalHeader)
}

func createDstFile(fileIndex int) (file *os.File, err os.Error) {
    fileName := fmt.Sprintf("%s%s%d%s", *dstDir, "/", fileIndex, ".pcap")
    return os.Open(fileName, os.O_CREAT | os.O_WRONLY, 0644)
}

func getPktSize(file *os.File, offset int64) (size int) {
    var buf [4]byte
    file.ReadAt(&buf, offset + 8)
    return 16 + int(buf[0]) + int(buf[1]) << 8 + int(buf[2]) << 16 + int(buf[3]) << 24 
}

func writeGlobalHeader(file *os.File) (n int, err os.Error) {
  return file.Write(&globalHeader)
}

func copyPacket(dst *os.File, src *os.File, offset int64, size int) (err os.Error) {
    buf := make([]byte, size)
    _, err = src.ReadAt(buf, offset)
    if err != nil && err != os.EOF {
        return err
    }
    dst.Write(buf) 
    return err
}

func main() {
	flag.Parse()

  log.Stdout(*toBeSplitted)
  log.Stdout(*dstDir)
  log.Stdout(*sizeOfSeg)

	srcFile, err := os.Open(*toBeSplitted, os.O_RDONLY, 0644)
  loadGlobalHeader(srcFile)

	os.MkdirAll(*dstDir, 0644)

  fileIndex := 1
  fileSize := int64(24)
  offset := int64(24)

  for err == nil {
      var dstFile *os.File
      dstFile, err = createDstFile(fileIndex)
      _, err = writeGlobalHeader(dstFile) 

      for fileSize < int64(*sizeOfSeg) * 1024 {
          pktSize := getPktSize(srcFile, offset)
          err = copyPacket(dstFile, srcFile, offset, pktSize)
          if err == nil {
              offset = offset + int64(pktSize)
              fileSize = fileSize + int64(pktSize)
          } else {
              break;
          }
      }
      dstFile.Close()
      fileSize = 24 
      fileIndex = fileIndex + 1
  }
  srcFile.Close()
}
