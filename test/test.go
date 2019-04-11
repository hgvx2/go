package main

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"proto_go"
	"tools/binaryBuf"
	"tools/log"
	"tools/server"
	"tools/server/serverDefine"
	"tools/server/tcpServer"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

/*
func daemon(nochdir, noclose int) int {
	var ret, ret2 uintptr
	var err syscall.Errno

	darwin := runtime.GOOS == "darwin"

	// already a daemon
	if syscall.Getppid() == 1 {
		return 0
	}

	// fork off the parent process
	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return -1
	}

	// failure
	if ret2 < 0 {
		os.Exit(-1)
	}

	// handle exception for darwin
	if darwin && ret2 == 1 {
		ret = 0
	}

	// if we got a good PID, then we call exit the parent process.
	if ret > 0 {
		os.Exit(0)
	}

	// Change the file mode mask
	_ = syscall.Umask(0)

	// create a new SID for the child process
	s_ret, s_errno := syscall.Setsid()
	if s_errno != nil {
		log.Printf("Error: syscall.Setsid errno: %d", s_errno)
	}
	if s_ret < 0 {
		return -1
	}

	if nochdir == 0 {
		os.Chdir("/")
	}

	if noclose == 0 {
		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if e == nil {
			fd := f.Fd()
			syscall.Dup2(int(fd), int(os.Stdin.Fd()))
			syscall.Dup2(int(fd), int(os.Stdout.Fd()))
			syscall.Dup2(int(fd), int(os.Stderr.Fd()))
		}
	}

	return 0
}
*/

// Product _
type Product struct {
	Name string       `json:"name"`
	ProductID int64   `json:"product_id"`
	Number int        `json:"number"`
	Price float64     `json:"price"`
	IsOnSale bool     `json:"is_on_sale"`
	TestArr []string  `json:"test_arr"`
}

var typeWork test.PhoneType = test.PhoneType_WORK
var typeHome test.PhoneType = test.PhoneType_WORK

func write() {

	//t1 := test.Phone{, proto.String("111111111"),}

	p1 := &test.Person{
		Id:   proto.Int32(1),
		Name: proto.String("小张"),
		Phones: []*test.Phone{
			{Type:&typeWork, Number:proto.String("1111111111"),},
			{Type:&typeHome, Number:proto.String("22222222222"),},
		},
	};
	p2 := &test.Person{
		Id:   proto.Int32(2),
		Name: proto.String("小往"),
		Phones: []*test.Phone{
			{Type:&typeWork, Number:proto.String("33333333333"),},
			{Type:&typeHome, Number:proto.String("444444444444"),},
		},
	};

	//创建地址簿
	book := &test.ContactBook{};
	book.Persons = append(book.Persons, p1);
	book.Persons = append(book.Persons, p2);

	//编码数据
	data, _ := proto.Marshal(book);
	//把数据写入文件
	ioutil.WriteFile("E:/GoWork/im/src/proto_go/test.txt", data, os.ModePerm);
}

func read()  {
	data, _ := ioutil.ReadFile("E:/GoWork/im/src/proto_go/test.txt")
	book := &test.ContactBook{};
	if nil != proto.Unmarshal(data, book){
		fmt.Println("error")
		return
	}
	fmt.Println(book)
}

func main() {

	log.InitInstance(log.Log_g_nInfo, "c:/","test.log")
	log.LogError("nihao %s %d", "123", 10)

	//write()
	//json.Marshal()
}


func testServer()  {
	if _, err := server.InitServerBase("0.0.0.0:18000", 5, 20, onServerBase); nil != err{
		fmt.Println("初始化失败")
		return
	}
	server.ServerRunning()
}

type DbWorker struct {
	//mysql data source name
	Dsn string
}
func testDB() {
	dbw := DbWorker{
		Dsn: "root:123456@tcp(127.0.0.1:3306)/test",
	}
	db, err := sql.Open("mysql",
		dbw.Dsn)
	if err != nil {
		fmt.Println("connect mysql err", err)
		return
	}
	fmt.Println("connect mysql o'k", err)
	defer db.Close()

	c, err := redis.Dial("tcp", "192.168.1.89:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	fmt.Println("Connect to redis ok")
	if _, err = c.Do("SET", "1", "2"); nil != err{
		fmt.Println(err)
		return
	}
	fmt.Println("Connect to redis ok")
	defer c.Close()
}



func onServerBase(msgType int, fd serverDefine.SOCKET_FD, info string, ptr binaryBuf.BinaryBufPtr) {
	switch msgType {
	case tcpServer.TCP_ACCEPT:
		fmt.Println("accept: ", "fd = ", fd, " msg: ", info)
	case tcpServer.TCP_DISCONNECT:
		fmt.Println("disscont: ", "fd = ", fd, " msg: ", info)
	case tcpServer.TCP_KEEPALIVE:
		fmt.Println("keepalive: ", "fd = ", fd, " msg: ", info)
	case tcpServer.TCP_PACK:
		//fmt.Println("pack: ", "fd = ", fd, " msg: ", info)
		//g_nums ++
		//binaryBuf.DestroyBinaryBuf(ptr)
		//fmt.Println("recv pack = ", g_nums)
		if err:= server.PushSendData(fd, ptr); nil != err{
			binaryBuf.DestroyBinaryBuf(ptr)
			fmt.Println(err)
			server.CloseSocket(fd)
			fmt.Println("push error*********")
		}
	default:
		fmt.Println("unknow msg ", msgType)
	}
}