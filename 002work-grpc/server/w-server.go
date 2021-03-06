package main

/**
 * @Author safoti
 * @Date Created in 2022/7/15
 * @Description  grpc 四种流程方式
 **/
import (
	"context"
	"go-grpc-demo/002work-grpc/pn"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
)

func main() {
  lis,err:= net.Listen("tcp",port)
  if err !=nil{
	  log.Fatalf("failed to listen: %v", err)
  }

   s:=  grpc.NewServer()
	pn.RegisterISrStreamServiceServer(s,&WoServer{})
	log.Println("开始监听，等待远程调用...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
// 常量：监听端口
const (
	port = ":50051"
)
//引入空的构造体
type WoServer struct {
	 pn.UnimplementedISrStreamServiceServer
}
//响应的方法

// 单项流式 ：单个请求，单个响应
func ( ws *WoServer) ISReqSingelrep( ctx context.Context,  req *pn.SServerRequest) (*pn.SClientResponse, error) {
	id :=  req.GetId()
	log.Fatal("1:收到请求：",id)
	//返回信息
	return &pn.SClientResponse{Id: id,Name: "1. name-" + strconv.Itoa(int(id))},nil
}


// 服务端流式 ：单个请求，集合响应
func (ws *WoServer) ISReqMultrep(  req *pn.SServerRequest,  ps pn.ISrStreamService_ISReqMultrepServer) error {
	// 取得请求参数
	id := req.GetId()

	// 打印请求参数
	log.Println("2. 收到请求:", id)

	// 返回多条记录
	for i := 0; i < 10; i++ {
		ps.Send(&pn.SClientResponse{Id: int32(i),Name: "2. name-" + strconv.Itoa(i)})
	}
  return nil
}
// 客户端流式 ：集合请求，单个响应
func (ws *WoServer) MUISReqMultrep(  ps  pn.ISrStreamService_MUISReqMultrepServer) error {
	var addVal int32 = 0
	for  {
		// 一次接受一条记录
		singleRequest, err := ps.Recv()
		// 不等于io.EOF表示这是条有效记录
		if err == io.EOF {
			log.Println("3. 客户端发送完毕")
			break
		} else if err != nil {
			log.Fatalln("3. 接收时发生异常", err)
			break
		} else {
				log.Println("3. 收到请求:", singleRequest.GetId())
			// 收完之后，执行SendAndClose返回数据并结束本次调用
			addVal += singleRequest.GetId()
		}
	}
	return ps.SendAndClose(&pn.SClientResponse{Id: addVal, Name: "3. name-" + strconv.Itoa(int(addVal))})
}

// 双向流式 ：集合请求，集合响应
func (ws *WoServer) SMUISReqMultrep(ps pn.ISrStreamService_SMUISReqMultrepServer) error {
	// 简单处理，对于收到的每一条记录都返回一个响应
	for {
		singleRequest, err := ps.Recv()

		// 不等于io.EOS表示这是条有效记录
		if err == io.EOF {
			log.Println("4. 接收完毕")
			return nil
		} else if err != nil {
			log.Fatalln("4. 接收时发生异常", err)
			return err
		} else {
			log.Println("4. 接收到数据", singleRequest.GetId())

			id := singleRequest.GetId()

			if sendErr := ps.Send(&pn.SClientResponse{Id: id, Name: "4. name-" + strconv.Itoa(int(id))}); sendErr != nil {
				log.Println("4. 返回数据异常数据", sendErr)
				return sendErr
			}
		}
	}
}


