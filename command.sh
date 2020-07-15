#!/bin/bash
export JENKINS_SLAVE_LABEL=work-node
export JENKINS_SLAVE_MODE=default
export JENKINS_SLAVE_NUM_EXECUTORS=1
export JENKINS_SLAVE_REMOTE_FS=/root/
export CONFIG="<slave><label>$JENKINS_SLAVE_LABEL</label><launcher class='hudson.slaves.JNLPLauncher' /><mode>$JENKINS_SLAVE_MODE</mode><numExecutors>$JENKINS_SLAVE_NUM_EXECUTORS</numExecutors><remoteFS>$JENKINS_SLAVE_REMOTE_FS</remoteFS></slave>"
export AGENT_NAME=`curl -s http://169.254.169.254/latest/meta-data/instance-id`
java -jar /root/jenkins-cli.jar -s http://$JENKINS_MASTER/ create-node $AGENT_NAME <<< "$CONFIG"
trap "java -jar /root/jenkins-cli.jar -s http://$JENKINS_MASTER/ delete-node $AGENT_NAME" EXIT
sleep 5
java -jar /root/agent.jar -jnlpUrl http://$JENKINS_MASTER/computer/$AGENT_NAME/slave-agent.jnlp
