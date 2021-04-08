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
      parallel {
        stage('Test Golang 1.13') {
          agent any
          tools {
            go 'golang-1.13'
          }
          environment {
            GO11MODULE = 'on'
          }
          steps {
            sh 'sudo apt-get install gcc -y'
            sh 'go get -u github.com/jstemmer/go-junit-report'
            sh 'go get .'
            sh 'go build'
            sh 'go test -v 2>&1 | go-junit-report > report.xml'
          }
          post {
              always {
                  junit 'report.xml'
              }
          }
        }
      }
      stage('Compile + Test') {
      parallel {
        stage('Test Golang 1.14') {
          agent any
          tools {
            go 'golang-1.14'
          }
          environment {
            GO11MODULE = 'on'
          }
          steps {
            sh 'sudo apt-get install gcc -y'
            sh 'go get -u github.com/jstemmer/go-junit-report'
            sh 'go get .'
            sh 'go build'
            sh 'go test -v 2>&1 | go-junit-report > report.xml'
          }
          post {
              always {
                  junit 'report.xml'
              }
          }
        }
      }
      stage('Test Golang 1.15') {
          agent any
          tools {
            go 'golang-1.15'
          }
          environment {
            GO11MODULE = 'on'
          }
          steps {
            sh 'sudo apt-get install gcc -y'
            sh 'go get -u github.com/jstemmer/go-junit-report'
            sh 'go get .'
            sh 'go build'
            sh 'go test -v 2>&1 | go-junit-report > report.xml'
          }
          post {
              always {
                  junit 'report.xml'
              }
          }
        }
      }
      stage('Test Golang 1.16') {
          agent any
          tools {
            go 'golang-1.16'
          }
          environment {
            GO11MODULE = 'on'
          }
          steps {
            sh 'sudo apt-get install gcc -y'
            sh 'go get -u github.com/jstemmer/go-junit-report'
            sh 'go get .'
            sh 'go build'
            sh 'go test -v 2>&1 | go-junit-report > report.xml'
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
