$(function() {
    // Init Ace editor
    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/monokai");
    editor.getSession().setMode("ace/mode/golang");

    var $stdout = $("#stdout");
    var $stderr = $("#stderr");
    var $compileError = $("#compile_error");

    var cleanResults = function() {
        $stdout.html("");
        $stderr.html("");
        $compileError.html("");
    };

    // The "Run" button
    $("#btn_run").on('click', function(){
        var postData = {
            version: 2,
            body: editor.getValue()
        };
        $.post("/compile/", postData, function(data) {
            cleanResults();
            var response = JSON.parse(data);
            var errMsg = response["Errors"];
            $compileError.html(errMsg);
            for (var i in response["Events"]) {
                var event = response["Events"][i];
                if (event["Kind"] == "stdout" && event["Delay"] == 0) {
                    $stdout.html(event["Message"]);
                }
                if (event["Kind"] == "stderr" && event["Delay"] == 0) {
                    $stderr.html(event["Message"]);
                }
            }
        });
    });
});
