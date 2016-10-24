$(function() {
    // Init Ace editor
    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/monokai");
    editor.getSession().setMode("ace/mode/golang");

    // Marddown
    marked.setOptions({
      renderer: new marked.Renderer(),
      gfm: true,
      tables: true,
      breaks: false,
      pedantic: false,
      sanitize: false,
      smartLists: true,
      smartypants: false
    });

    // Render description markdown
    var $description = $("#description");
    if (window.description) {
        $description.html(marked(window.description));
    }

    var $stdoutWrapper = $("#stdout_wrapper");
    var $stderrWrapper = $("#stderr_wrapper");
    var $compileErrorWrapper = $("#compile_error_wrapper");

    var $stdout = $("#stdout");
    var $stderr = $("#stderr");
    var $compileError = $("#compile_error");

    var showStdout = function(content) {
        $stdoutWrapper.fadeIn(400);
        $stdout.html(content);
    };

    var showStderr = function(content) {
        $stderrWrapper.fadeIn(400);
        $stderr.html(content);
    };

    var showCompileError = function(content) {
        $compileErrorWrapper.fadeIn(400);
        $compileError.html(content);
    };

    var cleanResults = function() {
        $stdout.html("");
        $stderr.html("");
        $compileError.html("");
    };

    var hideResults = function() {
        $stdoutWrapper.hide();
        $stderrWrapper.hide();
        $compileErrorWrapper.hide();
    };

    var equalsIgnoreLineEnding = function(s1, s2) {
        return s1.replace(/\r\n/g, '\n') == s2.replace(/\r\n/g, '\n');
    }

    var compileDone = function(data) {
        // Clean and hide all results
        hideResults();
        cleanResults();

        // Show stdout, stderr, compile error
        var response = JSON.parse(data);
        var errMsg = response["Errors"];
        var stdout = "";
        if (errMsg) {
            showCompileError(errMsg);
        } else {
            for (var i in response["Events"]) {
                var event = response["Events"][i];
                if (event["Kind"] == "stdout" && event["Delay"] == 0) {
                    stdout = event["Message"];
                    showStdout(stdout);
                }
                if (event["Kind"] == "stderr" && event["Delay"] == 0) {
                    showStderr(event["Message"]);
                }
            }

            // Check if stdout is as expected
            if (window.expectedStdout) {
                if (equalsIgnoreLineEnding(window.expectedStdout, stdout)) {
                    swal("Good job!", "You passed the test", "success")
                } else {
                    $.growl.warning({ message: "Not correct. :)" });
                }
            }
        }
    };

    var compileFailed = function(err) {
        $.growl.error({ message: err });
    };

    // The "Run" button
    $("#btn_run").on('click', function(){
        var postData = {
            version: 2,
            body: editor.getValue()
        };
        $.post("/compile/", postData).done(compileDone).fail(compileFailed);
    });
});
