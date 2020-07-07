package main

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "encoding/binary"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
)

type SegmentsData struct {
    Time   int32
    Offset int32
}

func main() {
    webMaskName := "./test.webmask"
    f, _ := os.Open(webMaskName)
    stat, err := f.Stat()
    if err != nil {
        panic(err)
    }
    size := stat.Size()

    outpath := filepath.Dir(webMaskName) + "/svg"
    CreateDir(outpath)
    log.Println("开始解压……", outpath)
    log.Println("开始解压，输出路径为:", outpath)

    //读取前16个字节获得tag，version，checkcode信息做校验
    buffer := make([]byte, 16)
    _, err = f.Read(buffer)
    if err != nil {
        panic(err)
    }

    tag := string(buffer[0:4])
    version := BytesToInt32(buffer[4:8])
    checkcode := BytesToInt32(buffer[8:9])

    segments := BytesToInt32(buffer[12:16])
    log.Println(fmt.Sprintf("TAG:%s,VERSION:%d,CHECKCODE:%d", tag, version, checkcode))
    log.Println("分段数为:", segments)

    segmentsData := []*SegmentsData{}
    for i := int32(0); i < segments; i++ {
        f.Read(buffer)
        //cur_offset, _ := f.Seek(0, os.SEEK_CUR)
        //fmt.Printf("current offset is %d\n", cur_offset)
        start := BytesToInt32(buffer[0:4])
        end := BytesToInt32(buffer[8:12])
        if start == 0 && end == 0 {
            time := BytesToInt32(buffer[4:8])
            offset := BytesToInt32(buffer[12:16])
            segmentsData = append(segmentsData, &SegmentsData{Time: time, Offset: offset})
        }
    }

    num := len(segmentsData)
    var length int32
    maskNum := 0
    for i := 0; i < num; i++ {
        offset := segmentsData[i].Offset

        if i < num-1 {
            length = segmentsData[i+1].Offset - offset
        } else {
            length = int32(size) - offset
        }

        buffer2 := make([]byte, length)
        f.Read(buffer2)
        r := UnFlate(buffer2)

        for j := 0; j < len(r); {
            offset := BytesToInt32(r[j : j+4])
            time := BytesToInt32(r[j+8 : j+16])
            svgBase64Data := string(r[j+12 : j+12+int(offset)])
            svgDatas := strings.Split(svgBase64Data, ";base64,")
            svgData, _ := base64.StdEncoding.DecodeString(svgDatas[1])
            filename := fmt.Sprintf("%s/%d.svg", outpath, time)
            ioutil.WriteFile(filename, svgData, 0777)
            maskNum++
            j = j + 12 + int(offset)
        }
    }
    log.Println("共包含蒙版数量:", maskNum)
    log.Println("数据解压完成")
}

func UnFlate(compressSrc []byte) []byte {
    b := bytes.NewReader(compressSrc)
    var out bytes.Buffer
    r, _ := gzip.NewReader(b)
    io.Copy(&out, r)
    return out.Bytes()
}

func BytesToInt32(bys []byte) (convInt int32) {
    bytebuff := bytes.NewBuffer(bys)
    binary.Read(bytebuff, binary.BigEndian, &convInt)
    return
}

func CreateDir(path string) {
    if !Exist(path) {
        os.MkdirAll(path, 0777)
    }
}

func Exist(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}
