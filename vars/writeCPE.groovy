import static com.sap.piper.Prerequisites.checkScript

void call(Map parameters = [:]) {
    final script = checkScript(this, parameters) ?: this
    // Dear reviewer remind me to fix this before merge
//     String piperGoPath = parameters.piperGoPath ?: './piper'
    String piperGoPath = './piper'
    Map cpe = script.commonPipelineEnvironment.getCPEMap(script)
    if (!cpe) {
        return
    }
    def writeCPECommand = """
cat <<EOF | ${piperGoPath} writeCPE
${groovy.json.JsonOutput.toJson(cpe)}
EOF
"""

    def output = script.sh(returnStdout: true, script: writeCPECommand)
    script.echo("Output is ${output}")

    output = script.sh(returnStdout: true, script: "ls -Ral .pipeline/commonPipelineEnvironment")
    script.echo("ls output is ${output}")
}
