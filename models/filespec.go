package models

import (
	"bufio"
	"fmt"
	"os"

	"github.com/saurbh7517/artifact/errorhandler"
)

type fileSpec struct {
	filename  string
	filesize  int
	blocksize int
	file      *os.File
	mapfile   *map[int][]byte
}

func NewFile(fileName string, fileSize int, blockSize int, file *os.File) *fileSpec {
	return &fileSpec{filename: fileName, filesize: fileSize, blocksize: blockSize, file: file}
}

func addByteToSlice(arr *[]byte, originalSlice *[]byte) {
	for i := 0; i < len(*arr); i++ {
		*originalSlice = append(*originalSlice, (*arr)[i])
	}
}

func writeToFile(fileIndex int, blockSlicePtr *[]byte) {
	// const outputFileLocation string = "E:/go_workspace/artifact"
	// fileName := outputFileLocation + "/" + "split_" + strconv.Itoa(fileIndex) + ".txt"
	// file, err := os.Create(fileName)
	// defer file.Close()
	// errorhandler.Check(err)
	// bufferedWriter := bufio.NewWriter(file)
	// bytesWritten, err := bufferedWriter.Write(*blockSlicePtr)
	// errorhandler.Check(err)
	// log.Println("The number of bytes written is ", bytesWritten)

	// bufferedWriter.Flush()

}

func (fs *fileSpec) ProcessFileByNewLine() *map[int][]byte {
	//This code will load a file in memory and split that into configurable size
	var numberOfBlocks int
	var fileByteSize int
	var byteCount int = 0
	const newLineCharacter byte = '\n'

	const blockByteSize int = 300 // block size of the file
	const lineLength int = 1000   // characters to be read when reading a file

	line := make([]byte, 0, lineLength)          // creating a slice in memory for reading a line
	blockSlice := make([]byte, 0, blockByteSize) //initializing a block size

	defer fs.file.Close()

	//Getting file details from memory, for detailed info check documentation
	fi, err := fs.file.Stat()
	errorhandler.Check(err)

	fileByteSize = int(fi.Size())

	/*Calculating the number of blocks*/
	numberOfBlocks = fileByteSize / blockByteSize
	if fileByteSize%blockByteSize != 0 {
		numberOfBlocks = numberOfBlocks + 1
	}
	/**********************************/
	mapper := make(map[int][]byte, numberOfBlocks) // initializing a map to store block size and its id
	scanner := bufio.NewScanner(fs.file)
	scanner.Split(bufio.ScanLines)

	eofStatus := scanner.Scan()
	count := 0
	// The following is the way to initialize a 2-D array in GO
	// blockSlice := make([][]byte, numberOfBlocks)
	// for i := range blockSlice {
	// 	blockSlice[i] = make([]byte, 0, blockByteSize)
	// }

	for eofStatus != false {
		line = scanner.Bytes()
		line = append(line, newLineCharacter)
		// fmt.Println(string(line))
		byteCount = byteCount + len(line)

		if byteCount == blockByteSize {
			//Write the contents of the block size to a file
			addByteToSlice(&line, &blockSlice)
			if count < numberOfBlocks {
				mapper[count+1] = blockSlice
				// writeToFile(count+1, &blockSlice)
				count++
			}
			// blockSlice = make([]byte, 0, blockByteSize)
			blockSlice = blockSlice[:0]
			byteCount = 0
		} else if byteCount > blockByteSize {
			if count < numberOfBlocks {
				mapper[count+1] = blockSlice
				// writeToFile(count+1, &blockSlice)
				count++
			}
			// blockSlice = make([]byte, 0, blockByteSize)
			blockSlice = blockSlice[:0] //doing this empty the array as the contents are copied by value in the map, so the same address space of the array can be re-used

			addByteToSlice(&line, &blockSlice)
			byteCount = len(line)
		} else {
			addByteToSlice(&line, &blockSlice)
		}

		eofStatus = scanner.Scan()
		// counter++
	}
	//Writing the remaining bytes to a file
	mapper[count+1] = blockSlice
	// writeToFile(count+1, &blockSlice)
	fmt.Println("The map is ", len(mapper))
	fs.mapfile = &mapper
	return fs.mapfile
}
