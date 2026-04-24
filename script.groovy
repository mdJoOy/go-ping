def testApp() {
        echo "Testing the program...."
    }
def buildApp() {
    when {
            expression {
                    BRANCH_NAME == 'main'
                }
        }
    echo "Building the program...."
    sh "docker build -t 192.168.0.45:8083/go-ping:1.0 ."

    }
def deployApp() {
    when {
            expression {
                    BRANCH_NAME == 'main'
                }
        }
    withCredentials([usernamePassword(credentialsId:'nexus-user-credentials', usernameVariable: 'USER', passwordVariable:'PWD')]){
        sh "echo '${PWD}' | docker login -u '${USER}' --password-stdint 192.168.0.45:8083"
        }
    sh "docker push 192.168.0.45:8083/go-ping:1.0"

    }
return this
