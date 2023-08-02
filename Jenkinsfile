pipeline {
  agent {
    node {
      label 'docker-agent-go'
    }

  }
  stages {
    stage('build') {
      steps {
        sh 'go mod tidy'
      }
    }

  }
}