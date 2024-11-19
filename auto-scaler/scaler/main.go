package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type appHandler struct{}

var Wait_goFunc sync.WaitGroup

func ageLimitScaleOut() {

	defer Wait_goFunc.Done()

	var deploymentName string

	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] AGE Scale-Out 조건이 설정되었습니다.")

	for {

		cmd_checkDeployment := exec.Command("bash", "-c", "kubectl get deploy -n microservice | grep 3m30s | awk '{print $1}' | head -1")
		output_deploymentName, err := cmd_checkDeployment.Output()
		if err != nil {
			fmt.Println(err)
		}

		str_deploymentName := string(output_deploymentName)
		deploymentName = str_deploymentName
		deploymentName = strings.ReplaceAll(deploymentName, "\n", "")

		if deploymentName == "" {
			continue
		} else {
			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "의 Age가 3m30s로 감지되었습니다.\n")
			time.Sleep(time.Second * 3)

			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + "Scale-Out을 진행합니다.\n")
			time.Sleep(time.Second * 3)

			retCmd := exec.Command("bash", "-c", "kubectl scale deployment --replicas=3 --namespace microservice "+deploymentName)
			_, err := retCmd.Output()
			if err != nil {
				fmt.Println(err)
			} else {
				time.Sleep(time.Second * 10)
				fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "이 정상적으로 Scale-Out 되었습니다.\n")
			}
		}
	}
}

func cpuLimitScaleOut() {

	defer Wait_goFunc.Done()

	var deploymentName string
	var hpaCount string
	var hpaName string
        
        hpabeforeWtachCmd := exec.Command("bash", "-c", "kubectl get hpa --namespace microservice")
	stdin, err := hpabeforeWtachCmd.StdinPipe()
        if err != nil {
		log.Panic(err)
        }
        go func() {
		defer stdin.Close()
                _, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
        }()
        beforeout, err := hpabeforeWtachCmd.CombinedOutput()
        if err != nil {
		fmt.Println("error", err)
        }
        fmt.Println("\n" + string(beforeout))

	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + "CPU Scale-Out 조건이 설정되었습니다.")

	for {

		CountCmd := exec.Command("bash", "-c", "kubectl get pods -n microservice | grep cpu | wc -l")
		output_CountCmd, err := CountCmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		str_Count := string(output_CountCmd)
		hpaCount = str_Count
		hpaCount = strings.ReplaceAll(hpaCount, "\n", "")

		int_Count, err := strconv.Atoi(hpaCount)
		if int(int_Count) >= 2 {

			cmd_checkHpa := exec.Command("bash", "-c", "kubectl get hpa -n microservice | awk '{print $1}'")
			output_hpaName, err := cmd_checkHpa.Output()
			if err != nil {
				fmt.Println(err)
			}

			str_hpaName := string(output_hpaName)
			hpaName = str_hpaName
			hpaName = strings.ReplaceAll(hpaName, "\n", "")

			if hpaName == "" {
				continue
			} else {
				time.Sleep(time.Second * 15)
				hpaWtachCmd := exec.Command("bash", "-c", "kubectl get hpa --namespace microservice")

				stdin, err := hpaWtachCmd.StdinPipe()
				if err != nil {
					log.Panic(err)
				}
				go func() {
					defer stdin.Close()
					_, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
				}()
				out, err := hpaWtachCmd.CombinedOutput()
				if err != nil {
					fmt.Println("error", err)
				}
				fmt.Println("\n" + string(out))

				time.Sleep(time.Second * 10)
				hpaDeleteCmd := exec.Command("bash", "-c", "kubectl delete hpa hpa-resource-cpu --namespace microservice")
				_, err = hpaDeleteCmd.Output()
				if err != nil {
					fmt.Println(err)
				}

				time.Sleep(time.Second * 5)
				cmd_checkDeployment := exec.Command("bash", "-c", "kubectl get deploy -n microservice | grep cpu | awk '{print $1}' | head -1")
				output_deploymentName, err := cmd_checkDeployment.Output()
				if err != nil {
					fmt.Println(err)
				}

				str_deploymentName := string(output_deploymentName)
				deploymentName = str_deploymentName
				deploymentName = strings.ReplaceAll(deploymentName, "\n", "")

				fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "이 정상적으로 Scale-Out 되었습니다.\n")
			}
		}
	}
}
func rpsScaleOut(){
	defer Wait_goFunc.Done()

	var deploymentName string
	var hpaCount string
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + "RPS Scale-Out 조건이 설정되었습니다.\n")
	var rpsFlag bool = false

	for !rpsFlag{
		
		CountCmd := exec.Command("bash", "-c", "kubectl get pods -n hpa | grep rps | wc -l")
		output_CountCmd, err := CountCmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		str_Count := string(output_CountCmd)
		hpaCount = str_Count
		hpaCount = strings.ReplaceAll(hpaCount, "\n", "")
		int_Count, err := strconv.Atoi(hpaCount)

		if int(int_Count) >= 2 {
			// HPA AutoScaleOut Success
			cmd_checkDeployment := exec.Command("bash", "-c", "kubectl get deploy -n hpa | grep rps | awk '{print $1}' | head -1")
				output_deploymentName, err := cmd_checkDeployment.Output()
				if err != nil {
					fmt.Println(err)
				}
			str_deploymentName := string(output_deploymentName)
			deploymentName = str_deploymentName
			deploymentName = strings.ReplaceAll(deploymentName, "\n", "")
			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "의 동작이 감지되었습니다.\n")
			time.Sleep(time.Second * 1)
			hpaWtachCmd := exec.Command("bash", "-c", "kubectl get hpa --namespace hpa")

				// fmt.Println("ok")
			stdin, err := hpaWtachCmd.StdinPipe()
			if err != nil {
				log.Panic(err)
			}
			go func() {
				defer stdin.Close()
				_, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
			}()
			out, err := hpaWtachCmd.CombinedOutput()
			if err != nil {
				fmt.Println("error", err)
			}
			fmt.Println("\n" + string(out))

			hpaDeleteCmd := exec.Command("bash", "-c", "kubectl delete hpa microservice-rps -n hpa")
			_, err = hpaDeleteCmd.Output()
			if err != nil {
					fmt.Println(err)
			}

			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "이 정상적으로 Scale-Out 되었습니다.\n")

			time.Sleep(time.Second * 10)
			podDeleteCmd := exec.Command("bash", "-c", "kubectl scale deployment microservice-rps -n hpa --replicas 1")
			_, err = podDeleteCmd.Output()
			if err != nil {
					fmt.Println(err)
			}
			rpsFlag = true	
		}else{
			continue
		}
	}

}
// Memory Usage based on HPA Auto ScaleOut Test Function
func memScaleOut(){
	defer Wait_goFunc.Done()
	var deploymentName string
	var hpaCount string

	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + "Memory Usage Scale-Out 조건이 설정되었습니다.\n")

	var memFlag bool = false

	for !memFlag{
		
		// namespace가 hpa로 설정된 파드 중 mem 문자열이 포함된 라인 수를 카운팅 후 반환하는 명령어
		// 즉, Pod name에 mem이 포함된 Pod의 개수를 반환
		// Default: 1, ScaleOut이 적용된 경우: 2 ~ n
		CountCmd := exec.Command("bash", "-c", "kubectl get pods -n hpa | grep mem | wc -l") 
		output_CountCmd, err := CountCmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		// String to integer
		str_Count := string(output_CountCmd)
		hpaCount = str_Count
		hpaCount = strings.ReplaceAll(hpaCount, "\n", "")
		int_Count, err := strconv.Atoi(hpaCount)

		// ScaleOut이 수행되어 Pod가 생성된 경우 수행되는 if문
		if int(int_Count) >= 2 {
			// HPA AutoScaleOut Success
			// hpa 네임스페이스 내의 mem이 포함된 파드의 정보 중 가장 앞의 Pod Name부분만 가져오는 명령어
			cmd_checkDeployment := exec.Command("bash", "-c", "kubectl get deploy -n hpa | grep mem | awk '{print $1}' | head -1")
				output_deploymentName, err := cmd_checkDeployment.Output()
				if err != nil {
					fmt.Println(err)
				}
			// 가져온 Pod Name을 Golang에서 사용하기 쉽도록 가공
			str_deploymentName := string(output_deploymentName)
			deploymentName = str_deploymentName
			deploymentName = strings.ReplaceAll(deploymentName, "\n", "")

			// ScaleOut 결과 Print
			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "의 동작이 감지되었습니다.\n")
			time.Sleep(time.Second * 1)

			// 현재 HPA Status 출력을 위한 커멘드
			hpaWtachCmd := exec.Command("bash", "-c", "kubectl get hpa --namespace hpa")
			// hpaWtachCmd의 표준 입력을 위한 파이프를 생성합니다.
			stdin, err := hpaWtachCmd.StdinPipe()
			if err != nil {
				log.Panic(err) // 오류가 발생하면 프로그램을 중단하고 오류 메시지를 출력합니다.
			}

			// goroutine을 생성하여 stdin에 데이터를 씁니다.
			go func() {
				defer stdin.Close() // 함수 종료 시 stdin을 닫아 명령어에 입력이 완료되었음을 알립니다.
				_, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input") // stdin에 문자열을 씁니다. 실제로 이 줄은 이 명령어에서는 의미가 없습니다.
			}()

			// 명령어를 실행하고 표준 출력 및 표준 에러 출력을 함께 가져옵니다.
			out, err := hpaWtachCmd.CombinedOutput()
			if err != nil {
				// 명령어 실행 중 오류가 발생한 경우 오류 메시지를 출력합니다.
				fmt.Println("error", err)
			}else{
				// 명령어의 출력 결과를 문자열로 변환하여 출력합니다.
				fmt.Println("\n" + string(out))
			}

			// ScaleOut이 완료된 HPA를 제거하기 위한 커멘드
			hpaDeleteCmd := exec.Command("bash", "-c", "kubectl delete hpa microservice-mem -n hpa")
			_, err = hpaDeleteCmd.Output()
			if err != nil {
					fmt.Println(err)
			}

			fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + deploymentName + "이 정상적으로 Scale-Out 되었습니다.\n")

			time.Sleep(time.Second * 10)
			podDeleteCmd := exec.Command("bash", "-c", "kubectl scale deployment microservice-mem -n hpa --replicas 1")
			_, err = podDeleteCmd.Output()
			if err != nil {
					fmt.Println(err)
			}
			
			// 무의미한 Loop를 막기위한 Flag 설정
			memFlag = true
		}else{
			continue
		}
	}
}

func main() {

	Wait_goFunc.Add(4)

	go ageLimitScaleOut()
	go cpuLimitScaleOut()
	go rpsScaleOut()
	go memScaleOut()

	Wait_goFunc.Wait()
}
