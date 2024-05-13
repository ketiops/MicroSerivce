FROM openjdk:8-jdk

ENV KAFKA_VERSION="3.6.0" SCALA_VERSION="2.12-3.6.0"
ENV KAFKA_HOME=/opt/kafka_"$SCALA_VERSION"-"$KAFKA_VERSION"

RUN apt-get update && apt-get install -y wget \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://downloads.apache.org/kafka/"$KAFKA_VERSION"/kafka_"$SCALA_VERSION".tgz -O /tmp/kafka_"$SCALA_VERSION".tgz \
    && tar xfz /tmp/kafka_"$SCALA_VERSION".tgz -C /opt \
    && rm /tmp/kafka_"$SCALA_VERSION".tgz

RUN wget https://golang.org/dl/go1.19.linux-amd64.tar.gz -O /tmp/go1.19.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz \
    && rm /tmp/go1.19.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

COPY . /opt/kafka_"$SCALA_VERSION"

WORKDIR /opt/kafka_"$SCALA_VERSION"
RUN go get github.com/go-sql-driver/mysql
EXPOSE 9092

ENV PATH=${PATH}:${KAFKA_HOME}/bin

CMD ["sh", "-c", "tail -f /dev/null"]