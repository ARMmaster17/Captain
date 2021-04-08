pipeline {
  agent any
  tools {
    go 'golang-1.16'
  }
  stages {
    stage('Compile') {
      steps {
        withCredentials([string(credentialsId: 'discord-server-webhook', variable: 'webhookURL')]) {
          discordSend link: env.BUILD_URL, title: 'Captain Build' + env.JOB_NAME, webhookURL: webhookURL, description: "Build started"
        }
        sh 'cd pkg/atc; go build'
      }
    }
    stage('Test') {
      steps {
        sh 'go test'
      }
    }
    stage('Notify') {
      steps {
        withCredentials([string(credentialsId: 'discord-server-webhook', variable: 'webhookURL')]) {
          discordSend link: env.BUILD_URL, title: 'Captain Build' + env.JOB_NAME, webhookURL: webhookURL, description: "Build SUCCESS"
        }
      }
    }
  }
}
