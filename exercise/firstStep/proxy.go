// 运行，跑命令：go run proxy.go之后直接用netcat测试 nc 127.0.0.1 1080

// socks5代理服务器实现
// socks5: 虽然是代理协议，但不能翻墙，明文传输
// 用途：企业内网安全性高，严格的防火墙策略，导致管理员访问资源也很麻烦，socks5相当于开了口子，让用户通过特定端口访问内部所有资源

// 1. client 与 socks5 server 协商
// 2. socks5 server 通过协商
// 3. client 发送请求给 socks5 server
// 4. socks5 server 和 host 建立tcp连接
// 5. host 返还响应给 socks5 server
// 6. socks5 server 返还状态给 client

// for each request:
// 7. client 发送数据给 socks5 server
// 8. socks5 server relay数据给 host
// 9. host 返还响应给 socks5 server
// 10. socks5 server 返还响应给 client
package firststep

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

// socks5协议
const socks5Ver = 0x05
const cmdBind = 0x01
const atypeIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	// 监听端口
	server, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		panic(err)
	}
	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go process(client)
	}
}

func process(conn net.Conn) {
	// 链接关掉
	defer conn.Close()
	// 缓冲流
	reader := bufio.NewReader(conn)
	// 认证
	err := auth(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
	}
	// 请求
	err = connect(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
	}
}

// socks5协商（认证阶段）
// 参数：reader 读取流，conn tcp链接
func auth(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
	// X’00’ NO AUTHENTICATION REQUIRED
	// X’02’ USERNAME/PASSWORD

	// 读取版本号--1字节
	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read ver failed:%w", err)
	}
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	// 读取methodSize--1字节
	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read methodSize failed:%w", err)
	}
	// 利用methodSize创建methods的缓冲区，用io.ReadFull填充
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method failed:%w", err)
	}
	// 截至目前，成功读取到三个字段，打印日志
	log.Println("ver", ver, "method", method)
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	// 按照协议，返回一个包，告诉浏览器，选择了哪种认证方式
	// socks5Ver: 协议版本 0x00: 无需认证
	_, err = conn.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	return nil
}

// socks5代理（请求阶段）
func connect(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05
	// CMD 0x01表示CONNECT请求
	// RSV 保留字段，值为0x00
	// ATYP(重点关注) 目标地址类型，DST.ADDR的数据对应这个字段的类型。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值
	// DST.PORT 目标端口，固定2个字节

	// 不一个字节一个字节读取，而是一次性读取4个字节(ver, cmd, rsv, atyp)
	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed:%w", err)
	}
	ver, cmd, atyp := buf[0], buf[1], buf[3]
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	if cmd != cmdBind {
		return fmt.Errorf("not supported cmd:%v", cmd)
	}
	addr := ""
	// 根据atype类型，读取不同的地址
	switch atyp {
	// IPV4 4个字节
	case atypeIPV4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("read atyp failed:%w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	// host 先读1个字节（host长度），再读host长度个字节，在用io.ReadFull填满
	case atypeHOST:
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostSize failed:%w", err)
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed:%w", err)
		}
		// host转换为字符串
		addr = string(host)
	// 暂时不实现，用的比较少
	case atypeIPV6:
		return errors.New("IPv6: no supported yet")
	default:
		return errors.New("invalid atyp")
	}
	// 读取端口 2个字节 还是用2个字节的缓冲区，用io.ReadFull填充
	// 这里复用字节长度为4的buf，然后切片取最后两个字节
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed:%w", err)
	}
	// 端口号转换为大端序，将会与该端口建立连接
	port := binary.BigEndian.Uint16(buf[:2])

	// 最后一步，真正和ip端口建立连接，双向传输数据
	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return fmt.Errorf("dial dst failed:%w", err)
	}
	// 函数结束，连接关闭
	defer dest.Close()

	log.Println("dial", addr, port)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT

	// 这里返回一个包，告诉浏览器，已经连接成功
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	// 建立浏览器和目标服务器的双向数据转换
	// 标准库io.Copy可以实现单向数据转换
	// Copy(dst io.Writer, src io.Reader)把src（只读流）的数据（用死循环）拷贝到dst（只写流）
	// 用WithCancel创建一个上下文，用于取消数据转换，在234行等待数据转换结束
	// 实现任何一个方向数据转换失败，都会取消另一个方向的数据转换
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 我们需要两个goroutine，一个用于从浏览器读取数据，一个用于向浏览器写入数据
	go func() {
		_, _ = io.Copy(dest, reader)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, dest)
		cancel()
	}()

	<-ctx.Done()
	return nil
}
