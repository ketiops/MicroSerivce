# **커스텀 메트릭를 활용한 마이크로서비스 단위 자동확장 기술**

커스텀 메트릭을 활용하여 마이크로서비스 단위의 서비스별 자동확장(Auto-Scale) 기술은 마이크로서비스의 워크로드를 리소스 사용량에 따라 적절한 서비스의 개수를 유지하기 위한 기술이다.

### **커스텀 메트릭을 통한 자동확장 기술 구현**

- 커스텀 메트릭을 타겟 메트릭으로 설정하여 HPA 배포
    - HPA를 통해 다양한 메트릭을 설정할 수 있다. 쿠버네티스의 기본 메트릭을 설정하여 CPU 사용량 등을 타겟 메트릭으로 설정할 수 있다. 그 외 커스텀된 메트릭을 사용하기 위해서는 프로메테우스 등 메트릭 수집 도구에서 확인이 가능한 메트릭만 설정이 가능하다. 본 문서에서는 프로메테우스를 통해 수집된 메트릭을 PromQL을 통해 커스텀한 메트릭을 타겟 메트릭으로 설정하여 자동확장을 수행한다.
1. 메모리 사용량을 기반으로 동작하는 HPA 리소스 정의 Yaml은 아래와 같다.
    
    ```yaml
    apiVersion: autoscaling/v2
    kind: HorizontalPodAutoscaler
    metadata:
      name: microservice-mem-hpa
      namespace: hpa
    spec:
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: microservice-mem
      minReplicas: 1
      maxReplicas: 10
      metrics:
      - type: Resource
        resource:
          name: memory
          target:
            type: Utilization
            averageUtilization: 70
      behavior:
        scaleUp:
          stabilizationWindowSeconds: 10
          policies:
            - type: Percent
              value: 10
              periodSeconds: 20
        scaleDown:
          stabilizationWindowSeconds: 30
          policies:
            - type: Percent
              value: 20
              periodSeconds: 60
    ```
    
2. 초당 서비스 요청량을 타겟 메트릭으로 동작하는 HPA 리소스 정의 파일은 아래와 같다.
    
    ```yaml
    apiVersion: autoscaling/v2
    kind: HorizontalPodAutoscaler
    metadata:
      name: microservice-mem-hpa
      namespace: hpa
    spec:
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: microservice-mem
      minReplicas: 1
      maxReplicas: 10
      metrics:
      - type: Resource
        resource:
          name: memory
          target:
            type: Utilization
            averageUtilization: 70
      behavior:
        scaleUp:
          stabilizationWindowSeconds: 10
          policies:
            - type: Percent
              value: 10
              periodSeconds: 20
        scaleDown:
          stabilizationWindowSeconds: 30
          policies:
            - type: Percent
              value: 20
              periodSeconds: 60
    ```
    
- **자동확장 검증 프로그램**
    - 자동확장 검증 프로그램은 Go언어로 작성되었으며, 실시간으로 현재 클러스터를 모니터링하여 구성한 HPA가 정상적으로 동작하는지 모니터링을 수행한다. 각 HPA에 대해 Go Routine으로 구성되어 있으며, 추후 업데이트에 용이하도록 설계되어 있다.
    - 실행방법은 아래와 같다.
        - /auto-scaler/scaler의 main.go 빌드 및 실행
        
        ```yaml
        go build -o scaler main.go
        
        ./scaler
        ```
        
    - 각 항목에 대한 자동확장 검증 과정은 아래와 같다.
        - CPU 사용량
            - 테스트용 Deployment를 실행할 경우 CPU 사용량 Brust 코드가 실행되는 이미지 배포
            - Brust에 의해 CPU사용량이 증가하여 기준 값 이상이 될 경우 스케일 아웃이 진행
        - Memory 사용량
            - 테스트용 Deployment는 요청에 따라 리소스 버스트를 수행하는 Deployment
            - Memory Brust 요청을 curl을 통해 수행하여 기준 값 이상인 경우 스케일 아웃 진행
            
            ```bash
            #!/bin/bash
            
            read -p 'Please enter stress module ip address: ' host_ip
            
            read -p 'Please enter the 1st arguments(duration): ' duration
            
            read -p 'Please enter the 2st arguments(mem_amount): ' mem_amount
            
            curl "http://$host_ip:5000/memory_stress?duration=$duration&mem_amount=$mem_amount" &
            ```
            
        - RPS
            - Request-Per-Second 수에 따른 자동 확장 검증
            - 반복적으로 curl을 수행하는 스크립트를 활용하여 RPS를 임의로 발생
            
            ```bash
            # serviceIP 입력:
            while true; do curl 10.102.208.101:8080/version; sleep 0.01; done
            ````
