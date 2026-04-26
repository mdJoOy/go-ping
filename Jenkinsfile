def gv
pipeline {
        agent any
        stages {
                stage(init) {
                        steps {
                            script {
                                echo 'loading the script'
                                gv = load 'script.groovy' 
                            }
                        }
                    }
                stage(test) {
                    steps {
                        script {
                            gv.testApp()
                        }
                    }
                    }
                stage(build) {
                    steps { 
                        script {
                            gv.buildApp()
                        }

                    }
                    }
                stage(deploy) {
                        steps {
                            script {
                                gv.deployApp()
                            }
                        }
                    }
            }
    }
//just testing the github webhook
