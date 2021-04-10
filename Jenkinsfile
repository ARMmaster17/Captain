pipeline {
  agent none
  stages {
    stage('Post Build Notification') {
      agent any
      steps {
        withCredentials([string(credentialsId: 'discord-server-webhook', variable: 'webhookURL')]) {
          discordSend link: env.BUILD_URL, title: 'Captain Build' + env.JOB_NAME, webhookURL: webhookURL, description: "Build started"
        }
      }
    }
    stage('Compile + Test') {
      matrix {
        agent any
        axes {
          axis {
            name 'GOLANG_VERSION'
            values 'golang-1.13', 'golang-1.14', 'golang-1.15', 'golang-1.16'
          }
          axis {
            name 'DATABASE_TYPE'
            values 'postgres', 'sqlite3-file' 'sqlite3-memory'
          }
        }

        stages {
          stage('Test ${GOLANG_VERSION}') {
            tools {
              go "${GOLANG_VERSION}"
            }
            environment {
              GO11MODULE = 'on'
              DATABASE_TYPE = "${DATABASE_TYPE}"
            }
            steps {
              sh 'go get -u github.com/jstemmer/go-junit-report'
              sh 'go get .'
              sh 'go build'
              sh 'mkdir /etc/captain'
              sh 'chmod 777 /etc/captain'
              sh 'go test -v 2>&1 | /home/administrator/go/bin/go-junit-report > report.xml'
              sh 'rm -rf /etc/captain'
            }
            post {
              always {
                junit 'report.xml'
              }
            }
          }
        }
      }
    }
  }
}
