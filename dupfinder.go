package main

import (
    "container/list"
    "crypto/md5"
    "bufio"
    "io"
    "os"
    "flag"
    "fmt"
    "log"
    "sync"
)

var wg sync.WaitGroup

type FileSum struct {
    file string
    checksum []byte
}

// check if filepath points to a directory
func isDirectory(absoluteFilePath string) (bool, error) {
    fileInfo, err := os.Stat(absoluteFilePath)
    if err != nil {
        return false, err
    }
    return fileInfo.Mode().IsDir(), err
}

func getChecksum(file *os.File) []byte {
    // calculate checksum
    hash := md5.New()
    _, err := io.Copy(hash, file)
    if err != nil {
        log.Fatal(err)
    }
    //log.Println(n, "bytes copied")
    return hash.Sum(nil)
}

func processFile(absoluteFilePath string, sums chan FileSum, worker chan int)  {
    // open file
    file, err := os.Open(absoluteFilePath)
    if err != nil {
        worker <- 1
        wg.Done()
    	log.Println(err)
        return
    }
    defer file.Close()
    // calculate checksum
    sum := getChecksum(file)

    // send md5 sum
    fs := FileSum{absoluteFilePath, sum}
    sums <- fs

    // release worker
    worker <- 1
}

func collectSums(sums chan FileSum, quit chan bool, collectedSums *map[string]*list.List)  {
    counter := 0
    //collectedSums := make(map[string]*list.List)
    for {
        select {
        case sum := <-sums:
            // read md5 sum
            counter += 1
            //log.Printf("%x\n", sum.checksum)
            //log.Println("collected sum", counter)

            // // // add sum to collectedSums
            // check whether checksum is already in map
            stringSum := fmt.Sprintf("%x", sum.checksum)
            fileList, ok := (*collectedSums)[stringSum]
            if ok {
                // checksum alraydy in map
                // add file to list
                fileList.PushBack(sum.file)
            } else {
                // new checksum found
                // add file to a new list
                fileList := list.New()
                fileList.PushBack(sum.file)
                // add list to collectedSums
                (*collectedSums)[stringSum] = fileList
            }
            // // // FINISH add sum to collectedSums

            // mark work as done
            wg.Done()
        case <-quit:
            log.Println("quit this shit")
            return
        }
    }
}

// read command line options
var workerCount int
func init() {
    flag.IntVar(&workerCount, "w", 4, "count of parallel md5sum workers")
}

func main() {
    flag.Parse()
    log.Println("Worker count", workerCount)

    // define worker count
    MAX_WORKER := workerCount
    worker := make(chan int, MAX_WORKER)
    // init worker pool
    for i := 0; i < MAX_WORKER; i++ {
        worker <- 1
    }

    // collector for checksums
    sums := make(chan FileSum)

    // start checksum collector
    quit := make(chan bool)
    collectedSums := make(map[string]*list.List)
    go collectSums(sums, quit, &collectedSums)

    // process all filenames read on stdin
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        // get next file
        file := scanner.Text()

        // skip directories
        isDirectory, err := isDirectory(file)
        if err != nil  {
            log.Println("ERROR", err)
            continue
        } else if (isDirectory) {
            //log.Println("Skip directory:", file)
            continue
        }
        // else process file
        <-worker
        wg.Add(1)
        go processFile(file, sums, worker)
    }

    // wait until all files sums are collected
    wg.Wait()
    quit <- true

    // print duplicates
    for checksum, files := range collectedSums {
        // check if list contains duplicates
        if files.Len() > 1 {
            fmt.Printf("Checksum %s:\n", checksum)
            for file := files.Front(); file != nil; file = file.Next() {
                fmt.Println("\t", file.Value)
            }
        }
    }
}
