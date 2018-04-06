const fs = require('fs');
const AdmZip = require('adm-zip');
const AWS = require('aws-sdk');
const zipper = new AdmZip();
console.log('Zipping ./dist directory...');
// Note: Does NOT work on Windows! See https://github.com/cthackers/adm-zip/pull/132

zipper.addLocalFolder('dist', '');
zipper.toBuffer(zipBuffer => {
    console.log('Successfully created zip buffer');
    const lambda = new AWS.Lambda({
        apiVersion: '2015-03-31',
        region: 'us-east-1',
        maxRetries: 3,
        sslEnabled: true,
        logger: console,
    });
    const params = {
        FunctionName: 'guten-snippet',
        Publish: false,
        ZipFile: zipBuffer,
    };
    lambda.updateFunctionCode(params, err => {
        if (err) {
            console.error(err);
            process.exit(3);
        }
    });
}, errString => {
    console.error(`Failed to create zip buffer: ${errString}`);
    process.exit(2);
}, filename => {
    console.log(`Zipping ${filename}...`);
});
