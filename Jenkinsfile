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
  }
}
