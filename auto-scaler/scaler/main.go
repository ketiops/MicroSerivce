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

func main() {

	Wait_goFunc.Add(3)

	go ageLimitScaleOut()
	go cpuLimitScaleOut()
	go rpsScaleOut()

	Wait_goFunc.Wait()
}
