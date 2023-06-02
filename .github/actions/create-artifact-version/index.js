const core = require('@actions/core');
const { exec } = require('child_process');

const GIT_COMMIT_DATE = "git log -1 --pretty='%ad' --date=format:'%Y%m%d.%H%M'";
const GIT_COMMIT_HASH = "git log -n 1 --pretty=format:'%h'";

exec(`echo "1.1_$(${GIT_COMMIT_DATE})_$(${GIT_COMMIT_HASH})"`, (error, stdout, stderr) => {
    if (error) {
        core.setFailed(error);
    }
    if (stderr) {
        core.setFailed(stderr);
    }
    core.setOutput("version", stdout);
});
