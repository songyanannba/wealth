package helper

import (
	"crypto/sha1"
	"fmt"
)

func Sha1(data []byte) string {
	//data := []byte("hello world")
	hash := sha1.Sum(data)

	//s := "string1"
	////产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	//h := sha1.New() // md5加密类似md5.New()
	////写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	//h.Write([]byte(s))
	////这个用来得到最终的散列值的字符切片。Sum 的参数可以用来对现有的字符切片追加额外的字节切片：一般不需要要。
	//bs := h.Sum(nil)
	////SHA1 值经常以 16 进制输出，使用%x 来将散列结果格式化为 16 进制字符串。
	//fmt.Printf("%x\n", bs)

	// 将哈希值转换成十六进制字符串
	hashHex := fmt.Sprintf("%x", hash)
	//fmt.Println(hashHex)
	return hashHex
}
