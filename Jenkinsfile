node {
  stage("Placeholder") {
    echo "Hello!"
  }
}

pipeline {
  agent any
  tools {
    go 'golang-1.16'
  }
  stages {
    stage('Compile') {
      steps {
        sh 'go build'
      }
    }
    stage('Test') {
      steps {
        sh 'go test'
      }
    }
    state('Notify') {
      steps {
        withCredentials([string(credentialsId: 'discord-server-webhook', variable: 'webhookURL')]) {
          discordSend link: env.BUILD_URL, title: 'Captain Build' + env.JOB_NAME, webhookURL: webhookURL
        }
      }
    }
  }
}
