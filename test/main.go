package main

import (
	"context"
	"github.com/bndr/gojenkins"
	"log"
)

func main() {
	// 初始化 Jenkins 客户端
	jenkins := gojenkins.CreateJenkins(nil, "http://192.168.20.16:8080", "zhiwei.zhang", "Rxnc9b7MS-Xo2Sse7UdJ")
	_, err := jenkins.Init(context.Background())
	if err != nil {
		log.Fatalf("Jenkins 初始化失败: %s", err)
	}

	//// 获取所有视图
	//views, err := jenkins.GetAllViews(context.Background())
	//if err != nil {
	//	log.Fatalf("获取 Jenkins 视图失败: %s", err)
	//}
	//
	//// 输出所有视图的名称
	//fmt.Println("Jenkins 视图:")
	//for _, view := range views {
	//	fmt.Println(view.GetName())
	//}

	//// 获取所有 jobs
	//jobs, err := jenkins.GetAllJobNames(context.Background())
	//if err != nil {
	//	log.Fatalf("获取 Jenkins jobs 失败: %s", err)
	//}
	//
	//fmt.Println("Jenkins jobs:")
	//for _, job := range jobs {
	//	fmt.Println(job.Name)
	//}

	//// 构建参数化任务
	//jobName := "example-job"
	//parameters := map[string]string{
	//	"param1": "value1",
	//	"param2": "value2",
	//}
	//
	//// 获取 job 对象
	//job, err := jenkins.GetJob(jobName)
	//if err != nil {
	//	log.Fatalf("获取 Jenkins job %s 失败: %s", jobName, err)
	//}
	//
	//// 参数化构建
	//queueID, err := job.InvokeSimple(parameters)
	//if err != nil {
	//	log.Fatalf("构建 Jenkins job %s 失败: %s", jobName, err)
	//}
	//
	//fmt.Printf("Job %s 开始构建，Queue ID: %d\n", jobName, queueID)
}
