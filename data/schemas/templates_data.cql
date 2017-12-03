// Copyright 2017 Verizon. All rights reserved.
// See provided LICENSE file for use of this source code.

//------------------------------------------------------------
// Single Device Templates
//------------------------------------------------------------

// 1. Number Card
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        77ea66da-9751-42c2-a6ec-5ede9d706d8b,
        'Single Device - Number Card',
        'Template used to generate a value card for a single ThingSpace device field.',
        'cell',
        '{"name":"Number Card","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","field":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBmaWVsZCA9IGNvbnRleHQuQXJnc1siZmllbGQiXQpsb2NhbCBzdWJmaWVsZCA9IGNvbnRleHQuQXJnc1sic3ViZmllbGQiXQpsb2NhbCBzZW1hbnRpYyA9IGNvbnRleHQuQXJnc1sic2VtYW50aWMiXQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWQpIHRoZW4KZXJyb3IoIkludmFsaWQgZGV2aWNlIGlkLiIpCmVuZAppZiBpc0VtcHR5KGZpZWxkKSB0aGVuCmVycm9yKCJJbnZhbGlkIGRldmljZSBmaWVsZC4iKQplbmQKbG9jYWwgdmFsdWUgPSBnZXREZXZpY2VGaWVsZFZhbHVlKGRldmljZUlkLCBmaWVsZCwgc3ViZmllbGQsIHNlbWFudGljKQpyZXR1cm4gZ2VuZXJhdGVWYWx1ZSh2YWx1ZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZCB2YWx1ZS4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRWYWx1ZShpZCwgZmllbGQsIHN1YmZpZWxkLCBzZW1hbnRpYykKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGQoaWQsIGZpZWxkLCBzdWJmaWVsZCkKaWYgZXJyIH49IG5pbCB0aGVuCmVycm9yKGVycikKZW5kCmlmIHJlc3AgPT0gbmlsIHRoZW4KZXJyb3IoIk5vIHJlY29yZHMgZm91bmQuIikKZW5kCmxvY2FsIHN1YmZpZWxkSW5kZXggPSAxCmlmIHN1YmZpZWxkIH49ICIiIHRoZW4KZm9yIGkgPSAxLCAjcmVzcC5jb2x1bW5zLCAxIGRvCmlmIHJlc3AuY29sdW1uc1tpXSA9PSBzdWJmaWVsZCB0aGVuCnN1YmZpZWxkSW5kZXggPSBpCmJyZWFrCmVuZAplbmQKZW5kCmxvY2FsIHZhbHVlID0gewp2YWx1ZSA9IHJlc3Aucm93c1sxXVtzdWJmaWVsZEluZGV4XSwKc2VtYW50aWMgPSBzZW1hbnRpYwp9CnJldHVybiB2YWx1ZQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgdmFsdWUuCmZ1bmN0aW9uIGdlbmVyYXRlVmFsdWUodmFsdWUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnZhbHVlKHZhbHVlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGVtcHR5IHN0cmluZy4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":null}}}',
        2,
        'U+RLrj9vsXPRMzoY+qVQm0qUNgg='
    );

// 2. Bar Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        f694361c-5b1b-4434-92a2-18027b2eba33,
        'Single Device - Bar Graph',
        'Template used to generate a Bar Graph for one or more fields on a ThingSpace device.',
        'cell',
        '{"name":"Bar Graph","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","fields":"","labels":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBkZXZpY2VGaWVsZHMgPSBjb250ZXh0LkFyZ3NbImZpZWxkcyJdCmxvY2FsIGZpZWxkTGFiZWxzID0gY29udGV4dC5BcmdzWyJsYWJlbHMiXQpsb2NhbCBsaW1pdCA9IHRvbnVtYmVyKGNvbnRleHQuQXJnc1sibGltaXQiXSkKLS0gVmVyaWZ5IHJlcXVpcmVkIGlucHV0cwppZiBpc0VtcHR5KGRldmljZUlkKSB0aGVuCmVycm9yKCJkZXZpY2VJZCBpcyBlbXB0eS4iKQplbmQKaWYgaXNFbXB0eShkZXZpY2VGaWVsZHMpIHRoZW4KZXJyb3IoImZpZWxkcyBpcyBlbXB0eS4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGZpZWxkcy4KbG9jYWwgZmllbGRzID0gc3BsaXQoZGV2aWNlRmllbGRzLCAiLCIpCmxvY2FsIGxhYmVscyA9IHNwbGl0KGZpZWxkTGFiZWxzLCAiLCIgKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZShzKSBmb3IgZGV2aWNlCmZvciBpID0gMSwgI2ZpZWxkcywgMSBkbwpmaWVsZCA9IHNwbGl0KGZpZWxkc1tpXSwgIi8iKQpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhkZXZpY2VJZCwgZmllbGRbMV0sIGZpZWxkWzJdLCBsaW1pdCkKbG9jYWwgY29sdW1uTmFtZSA9IGRldmljZVRhYmxlLmNvbHVtbnNbMV0KaWYgbGFiZWxzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBsYWJlbHNbaV0KZW5kCmRhdGFUYWJsZS5jb2x1bW5zW2ldID0gY29sdW1uTmFtZQpmb3IgaiA9IDEsICNkZXZpY2VUYWJsZS5yb3dzLCAxIGRvCmlmIGRhdGFUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRhdGFUYWJsZS5yb3dzW2pdID0ge30KZW5kCnRhYmxlLmluc2VydChkYXRhVGFibGUucm93c1tqXSwgZGV2aWNlVGFibGUucm93c1tqXVsxXSkKZW5kCmVuZApyZXR1cm4gZ2VuZXJhdGVUYWJsZShkYXRhVGFibGUpCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZXQgVGhpbmdTcGFjZSBkZXZpY2UgZmllbGQuCmZ1bmN0aW9uIGdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgc3ViZmllbGQgPT0gbmlsIHRoZW4Kc3ViZmllbGQgPSAiIgplbmQKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQppZiBlcnIgfj0gbmlsIHRoZW4KZXJyb3IoZXJyKQplbmQKaWYgcmVzcCA9PSBuaWwgdGhlbgplcnJvcigibm8gcmVjb3JkcyBmb3VuZCIpCmVuZApyZXR1cm4gcmVzcAplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgTm9ydGhTdGFyIHRhYmxlLgpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQpsb2NhbCBvdXQsIGVyciA9IG91dHB1dC50YWJsZSh0YWJsZSkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBzcGxpdCBkZWxpbWl0ZWQgc3RyaW5ncyBpbnRvIGx1YSB0YWJsZS4KZnVuY3Rpb24gc3BsaXQoc3RyLCBkZWxpbSkKaWYgc3RyaW5nLmZpbmQoc3RyLCBkZWxpbSkgPT0gbmlsIHRoZW4KcmV0dXJuIHsgc3RyIH0KZW5kCmxvY2FsIGFycmF5ID0ge30KbG9jYWwgaW5kZXggPSAxCmZvciB2YWx1ZSBpbiBzdHJpbmcuZ21hdGNoKHN0ciwgIihbXiIgLi4gZGVsaW0gLi4gIl0rKSIpIGRvCmFycmF5W2luZGV4XSA9IHZhbHVlCmluZGV4ID0gaW5kZXggKyAxCmVuZApyZXR1cm4gYXJyYXkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGlmIHN0cmluZyBpcyBlbXB0eS4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"discreteBarChart"}}}}',
        2,
        'haC9mQ7vVhIv8NqgS7qt2jA2Les='
    );

// 3. MultiBar Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        09c74f9a-787f-41a3-a07b-bb4983bb29a0,
        'Single Device - MultiBar Graph',
        'Template used to generate a MultiBar Graph for one or more fields on a ThingSpace device.',
        'cell',
        '{"name":"MultiBar Graph","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","fields":"","labels":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBkZXZpY2VGaWVsZHMgPSBjb250ZXh0LkFyZ3NbImZpZWxkcyJdCmxvY2FsIGZpZWxkTGFiZWxzID0gY29udGV4dC5BcmdzWyJsYWJlbHMiXQpsb2NhbCBsaW1pdCA9IHRvbnVtYmVyKGNvbnRleHQuQXJnc1sibGltaXQiXSkKLS0gVmVyaWZ5IHJlcXVpcmVkIGlucHV0cwppZiBpc0VtcHR5KGRldmljZUlkKSB0aGVuCmVycm9yKCJkZXZpY2VJZCBpcyBlbXB0eS4iKQplbmQKaWYgaXNFbXB0eShkZXZpY2VGaWVsZHMpIHRoZW4KZXJyb3IoImZpZWxkcyBpcyBlbXB0eS4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGZpZWxkcy4KbG9jYWwgZmllbGRzID0gc3BsaXQoZGV2aWNlRmllbGRzLCAiLCIpCmxvY2FsIGxhYmVscyA9IHNwbGl0KGZpZWxkTGFiZWxzLCAiLCIgKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZShzKSBmb3IgZGV2aWNlCmZvciBpID0gMSwgI2ZpZWxkcywgMSBkbwpmaWVsZCA9IHNwbGl0KGZpZWxkc1tpXSwgIi8iKQpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhkZXZpY2VJZCwgZmllbGRbMV0sIGZpZWxkWzJdLCBsaW1pdCkKbG9jYWwgY29sdW1uTmFtZSA9IGRldmljZVRhYmxlLmNvbHVtbnNbMV0KaWYgbGFiZWxzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBsYWJlbHNbaV0KZW5kCmRhdGFUYWJsZS5jb2x1bW5zW2ldID0gY29sdW1uTmFtZQpmb3IgaiA9IDEsICNkZXZpY2VUYWJsZS5yb3dzLCAxIGRvCmlmIGRhdGFUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRhdGFUYWJsZS5yb3dzW2pdID0ge30KZW5kCnRhYmxlLmluc2VydChkYXRhVGFibGUucm93c1tqXSwgZGV2aWNlVGFibGUucm93c1tqXVsxXSkKZW5kCmVuZApyZXR1cm4gZ2VuZXJhdGVUYWJsZShkYXRhVGFibGUpCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZXQgVGhpbmdTcGFjZSBkZXZpY2UgZmllbGQuCmZ1bmN0aW9uIGdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgc3ViZmllbGQgPT0gbmlsIHRoZW4Kc3ViZmllbGQgPSAiIgplbmQKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQppZiBlcnIgfj0gbmlsIHRoZW4KZXJyb3IoZXJyKQplbmQKaWYgcmVzcCA9PSBuaWwgdGhlbgplcnJvcigibm8gcmVjb3JkcyBmb3VuZCIpCmVuZApyZXR1cm4gcmVzcAplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgTm9ydGhTdGFyIHRhYmxlLgpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQpsb2NhbCBvdXQsIGVyciA9IG91dHB1dC50YWJsZSh0YWJsZSkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBzcGxpdCBkZWxpbWl0ZWQgc3RyaW5ncyBpbnRvIGx1YSB0YWJsZS4KZnVuY3Rpb24gc3BsaXQoc3RyLCBkZWxpbSkKaWYgc3RyaW5nLmZpbmQoc3RyLCBkZWxpbSkgPT0gbmlsIHRoZW4KcmV0dXJuIHsgc3RyIH0KZW5kCmxvY2FsIGFycmF5ID0ge30KbG9jYWwgaW5kZXggPSAxCmZvciB2YWx1ZSBpbiBzdHJpbmcuZ21hdGNoKHN0ciwgIihbXiIgLi4gZGVsaW0gLi4gIl0rKSIpIGRvCmFycmF5W2luZGV4XSA9IHZhbHVlCmluZGV4ID0gaW5kZXggKyAxCmVuZApyZXR1cm4gYXJyYXkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGlmIHN0cmluZyBpcyBlbXB0eS4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"multiBarChart"}}}}',
        2,
        'haC9mQ7vVhIv8NqgS7qt2jA2Les='

    );

// 4. Line Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        77dd93c0-32a1-4c9a-822e-4b9fa69a9ed3,
        'Single Device - Line Graph',
        'Template used to generate a Line Graph for one or more fields on a ThingSpace device.',
        'cell',
        '{"name":"Line Graph","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","fields":"","labels":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBkZXZpY2VGaWVsZHMgPSBjb250ZXh0LkFyZ3NbImZpZWxkcyJdCmxvY2FsIGZpZWxkTGFiZWxzID0gY29udGV4dC5BcmdzWyJsYWJlbHMiXQpsb2NhbCBsaW1pdCA9IHRvbnVtYmVyKGNvbnRleHQuQXJnc1sibGltaXQiXSkKLS0gVmVyaWZ5IHJlcXVpcmVkIGlucHV0cwppZiBpc0VtcHR5KGRldmljZUlkKSB0aGVuCmVycm9yKCJkZXZpY2VJZCBpcyBlbXB0eS4iKQplbmQKaWYgaXNFbXB0eShkZXZpY2VGaWVsZHMpIHRoZW4KZXJyb3IoImZpZWxkcyBpcyBlbXB0eS4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGZpZWxkcy4KbG9jYWwgZmllbGRzID0gc3BsaXQoZGV2aWNlRmllbGRzLCAiLCIpCmxvY2FsIGxhYmVscyA9IHNwbGl0KGZpZWxkTGFiZWxzLCAiLCIgKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZShzKSBmb3IgZGV2aWNlCmZvciBpID0gMSwgI2ZpZWxkcywgMSBkbwpmaWVsZCA9IHNwbGl0KGZpZWxkc1tpXSwgIi8iKQpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhkZXZpY2VJZCwgZmllbGRbMV0sIGZpZWxkWzJdLCBsaW1pdCkKbG9jYWwgY29sdW1uTmFtZSA9IGRldmljZVRhYmxlLmNvbHVtbnNbMV0KaWYgbGFiZWxzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBsYWJlbHNbaV0KZW5kCmRhdGFUYWJsZS5jb2x1bW5zW2ldID0gY29sdW1uTmFtZQpmb3IgaiA9IDEsICNkZXZpY2VUYWJsZS5yb3dzLCAxIGRvCmlmIGRhdGFUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRhdGFUYWJsZS5yb3dzW2pdID0ge30KZW5kCnRhYmxlLmluc2VydChkYXRhVGFibGUucm93c1tqXSwgZGV2aWNlVGFibGUucm93c1tqXVsxXSkKZW5kCmVuZApyZXR1cm4gZ2VuZXJhdGVUYWJsZShkYXRhVGFibGUpCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZXQgVGhpbmdTcGFjZSBkZXZpY2UgZmllbGQuCmZ1bmN0aW9uIGdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgc3ViZmllbGQgPT0gbmlsIHRoZW4Kc3ViZmllbGQgPSAiIgplbmQKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQppZiBlcnIgfj0gbmlsIHRoZW4KZXJyb3IoZXJyKQplbmQKaWYgcmVzcCA9PSBuaWwgdGhlbgplcnJvcigibm8gcmVjb3JkcyBmb3VuZCIpCmVuZApyZXR1cm4gcmVzcAplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgTm9ydGhTdGFyIHRhYmxlLgpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQpsb2NhbCBvdXQsIGVyciA9IG91dHB1dC50YWJsZSh0YWJsZSkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBzcGxpdCBkZWxpbWl0ZWQgc3RyaW5ncyBpbnRvIGx1YSB0YWJsZS4KZnVuY3Rpb24gc3BsaXQoc3RyLCBkZWxpbSkKaWYgc3RyaW5nLmZpbmQoc3RyLCBkZWxpbSkgPT0gbmlsIHRoZW4KcmV0dXJuIHsgc3RyIH0KZW5kCmxvY2FsIGFycmF5ID0ge30KbG9jYWwgaW5kZXggPSAxCmZvciB2YWx1ZSBpbiBzdHJpbmcuZ21hdGNoKHN0ciwgIihbXiIgLi4gZGVsaW0gLi4gIl0rKSIpIGRvCmFycmF5W2luZGV4XSA9IHZhbHVlCmluZGV4ID0gaW5kZXggKyAxCmVuZApyZXR1cm4gYXJyYXkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGlmIHN0cmluZyBpcyBlbXB0eS4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"lineChart"}}}}',
        2,
        'haC9mQ7vVhIv8NqgS7qt2jA2Les='

    );

// 5. Pie Chart Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        541b8f82-7b11-4ee0-b5cd-baeacaaf812f,
        'Single Device - Pie Chart Graph',
        'Template used to generate a Pie Chart Graph for one or more fields on a ThingSpace device.',
        'cell',
        '{"name":"Pie Chart Graph","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","fields":"","labels":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBkZXZpY2VGaWVsZHMgPSBjb250ZXh0LkFyZ3NbImZpZWxkcyJdCmxvY2FsIGZpZWxkTGFiZWxzID0gY29udGV4dC5BcmdzWyJsYWJlbHMiXQpsb2NhbCBsaW1pdCA9IHRvbnVtYmVyKGNvbnRleHQuQXJnc1sibGltaXQiXSkKLS0gVmVyaWZ5IHJlcXVpcmVkIGlucHV0cwppZiBpc0VtcHR5KGRldmljZUlkKSB0aGVuCmVycm9yKCJkZXZpY2VJZCBpcyBlbXB0eS4iKQplbmQKaWYgaXNFbXB0eShkZXZpY2VGaWVsZHMpIHRoZW4KZXJyb3IoImZpZWxkcyBpcyBlbXB0eS4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGZpZWxkcy4KbG9jYWwgZmllbGRzID0gc3BsaXQoZGV2aWNlRmllbGRzLCAiLCIpCmxvY2FsIGxhYmVscyA9IHNwbGl0KGZpZWxkTGFiZWxzLCAiLCIgKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZShzKSBmb3IgZGV2aWNlCmZvciBpID0gMSwgI2ZpZWxkcywgMSBkbwpmaWVsZCA9IHNwbGl0KGZpZWxkc1tpXSwgIi8iKQpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhkZXZpY2VJZCwgZmllbGRbMV0sIGZpZWxkWzJdLCBsaW1pdCkKbG9jYWwgY29sdW1uTmFtZSA9IGRldmljZVRhYmxlLmNvbHVtbnNbMV0KaWYgbGFiZWxzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBsYWJlbHNbaV0KZW5kCmRhdGFUYWJsZS5jb2x1bW5zW2ldID0gY29sdW1uTmFtZQpmb3IgaiA9IDEsICNkZXZpY2VUYWJsZS5yb3dzLCAxIGRvCmlmIGRhdGFUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRhdGFUYWJsZS5yb3dzW2pdID0ge30KZW5kCnRhYmxlLmluc2VydChkYXRhVGFibGUucm93c1tqXSwgZGV2aWNlVGFibGUucm93c1tqXVsxXSkKZW5kCmVuZApyZXR1cm4gZ2VuZXJhdGVUYWJsZShkYXRhVGFibGUpCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZXQgVGhpbmdTcGFjZSBkZXZpY2UgZmllbGQuCmZ1bmN0aW9uIGdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgc3ViZmllbGQgPT0gbmlsIHRoZW4Kc3ViZmllbGQgPSAiIgplbmQKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQppZiBlcnIgfj0gbmlsIHRoZW4KZXJyb3IoZXJyKQplbmQKaWYgcmVzcCA9PSBuaWwgdGhlbgplcnJvcigibm8gcmVjb3JkcyBmb3VuZCIpCmVuZApyZXR1cm4gcmVzcAplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgTm9ydGhTdGFyIHRhYmxlLgpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQpsb2NhbCBvdXQsIGVyciA9IG91dHB1dC50YWJsZSh0YWJsZSkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBzcGxpdCBkZWxpbWl0ZWQgc3RyaW5ncyBpbnRvIGx1YSB0YWJsZS4KZnVuY3Rpb24gc3BsaXQoc3RyLCBkZWxpbSkKaWYgc3RyaW5nLmZpbmQoc3RyLCBkZWxpbSkgPT0gbmlsIHRoZW4KcmV0dXJuIHsgc3RyIH0KZW5kCmxvY2FsIGFycmF5ID0ge30KbG9jYWwgaW5kZXggPSAxCmZvciB2YWx1ZSBpbiBzdHJpbmcuZ21hdGNoKHN0ciwgIihbXiIgLi4gZGVsaW0gLi4gIl0rKSIpIGRvCmFycmF5W2luZGV4XSA9IHZhbHVlCmluZGV4ID0gaW5kZXggKyAxCmVuZApyZXR1cm4gYXJyYXkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGlmIHN0cmluZyBpcyBlbXB0eS4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"pieChart"}}}}',
        2,
        'haC9mQ7vVhIv8NqgS7qt2jA2Les='
    );

// 6. Data Table
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        c7d861b0-7085-4cd2-a023-89851dba16af,
        'Single Device - Data Table',
        'Template used to generate a Data Table for one or more fields on a ThingSpace device',
        'cell',
        '{"name":"Data Table","input":{"type":"Code","language":"lua","arguments":{"deviceId":"","fields":"","labels":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZCA9IGNvbnRleHQuQXJnc1siZGV2aWNlSWQiXQpsb2NhbCBkZXZpY2VGaWVsZHMgPSBjb250ZXh0LkFyZ3NbImZpZWxkcyJdCmxvY2FsIGZpZWxkTGFiZWxzID0gY29udGV4dC5BcmdzWyJsYWJlbHMiXQpsb2NhbCBsaW1pdCA9IHRvbnVtYmVyKGNvbnRleHQuQXJnc1sibGltaXQiXSkKLS0gVmVyaWZ5IHJlcXVpcmVkIGlucHV0cwppZiBpc0VtcHR5KGRldmljZUlkKSB0aGVuCmVycm9yKCJkZXZpY2VJZCBpcyBlbXB0eS4iKQplbmQKaWYgaXNFbXB0eShkZXZpY2VGaWVsZHMpIHRoZW4KZXJyb3IoImZpZWxkcyBpcyBlbXB0eS4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGZpZWxkcy4KbG9jYWwgZmllbGRzID0gc3BsaXQoZGV2aWNlRmllbGRzLCAiLCIpCmxvY2FsIGxhYmVscyA9IHNwbGl0KGZpZWxkTGFiZWxzLCAiLCIgKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZShzKSBmb3IgZGV2aWNlCmZvciBpID0gMSwgI2ZpZWxkcywgMSBkbwpmaWVsZCA9IHNwbGl0KGZpZWxkc1tpXSwgIi8iKQpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhkZXZpY2VJZCwgZmllbGRbMV0sIGZpZWxkWzJdLCBsaW1pdCkKbG9jYWwgY29sdW1uTmFtZSA9IGRldmljZVRhYmxlLmNvbHVtbnNbMV0KaWYgbGFiZWxzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBsYWJlbHNbaV0KZW5kCmRhdGFUYWJsZS5jb2x1bW5zW2ldID0gY29sdW1uTmFtZQpmb3IgaiA9IDEsICNkZXZpY2VUYWJsZS5yb3dzLCAxIGRvCmlmIGRhdGFUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRhdGFUYWJsZS5yb3dzW2pdID0ge30KZW5kCnRhYmxlLmluc2VydChkYXRhVGFibGUucm93c1tqXSwgZGV2aWNlVGFibGUucm93c1tqXVsxXSkKZW5kCmVuZApyZXR1cm4gZ2VuZXJhdGVUYWJsZShkYXRhVGFibGUpCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZXQgVGhpbmdTcGFjZSBkZXZpY2UgZmllbGQuCmZ1bmN0aW9uIGdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgc3ViZmllbGQgPT0gbmlsIHRoZW4Kc3ViZmllbGQgPSAiIgplbmQKbG9jYWwgcmVzcCwgZXJyID0gdHMuZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQppZiBlcnIgfj0gbmlsIHRoZW4KZXJyb3IoZXJyKQplbmQKaWYgcmVzcCA9PSBuaWwgdGhlbgplcnJvcigibm8gcmVjb3JkcyBmb3VuZCIpCmVuZApyZXR1cm4gcmVzcAplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2VuZXJhdGUgTm9ydGhTdGFyIHRhYmxlLgpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQpsb2NhbCBvdXQsIGVyciA9IG91dHB1dC50YWJsZSh0YWJsZSkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBzcGxpdCBkZWxpbWl0ZWQgc3RyaW5ncyBpbnRvIGx1YSB0YWJsZS4KZnVuY3Rpb24gc3BsaXQoc3RyLCBkZWxpbSkKaWYgc3RyaW5nLmZpbmQoc3RyLCBkZWxpbSkgPT0gbmlsIHRoZW4KcmV0dXJuIHsgc3RyIH0KZW5kCmxvY2FsIGFycmF5ID0ge30KbG9jYWwgaW5kZXggPSAxCmZvciB2YWx1ZSBpbiBzdHJpbmcuZ21hdGNoKHN0ciwgIihbXiIgLi4gZGVsaW0gLi4gIl0rKSIpIGRvCmFycmF5W2luZGV4XSA9IHZhbHVlCmluZGV4ID0gaW5kZXggKyAxCmVuZApyZXR1cm4gYXJyYXkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGNoZWNrIGlmIHN0cmluZyBpcyBlbXB0eS4KZnVuY3Rpb24gaXNFbXB0eShzdHIpCnJldHVybiBzdHIgPT0gbmlsIG9yIHN0ciA9PSAiIgplbmQ=","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"table"}}}}',
        2,
        'haC9mQ7vVhIv8NqgS7qt2jA2Les='
    );

// 7. Time Series Chart

// 8. Gauge


//------------------------------------------------------------
// Multiple Devices Templates
//------------------------------------------------------------

// 1. Bar Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        3156e252-9937-4c25-bd59-fb8aaaf9350e,
        'Multiple Devices - Bar Graph',
        'Template used to generate a Bar Graph for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"Bar Graph","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIHN1YmZpZWxkID0gY29udGV4dC5BcmdzWyJzdWJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJkZXZpY2VJZHMgaXMgZW1wdHkuIikKZW5kCmlmIGlzRW1wdHkoZmllbGQpIHRoZW4KZXJyb3IoImZpZWxkIGlzIGVtcHR5LiIpCmVuZAotLSBQYXJzZSBkZXZpY2UgaWRzLgpsb2NhbCBpZHMgPSBzcGxpdChkZXZpY2VJZHMsICIsIikKbG9jYWwgbmFtZXMgPSBzcGxpdChkZXZpY2VOYW1lcywgIiwiKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZSBmb3IgZGV2aWNlKHMpCmZvciBpID0gMSwgI2lkcywgMSBkbwpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhpZHNbaV0sIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIGNvbHVtbk5hbWUgPSBkZXZpY2VUYWJsZS5jb2x1bW5zWzFdCmlmIG5hbWVzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBuYW1lc1tpXQplbmQKZGF0YVRhYmxlLmNvbHVtbnNbaV0gPSBjb2x1bW5OYW1lCmZvciBqID0gMSwgI2RldmljZVRhYmxlLnJvd3MsIDEgZG8KaWYgZGF0YVRhYmxlLnJvd3Nbal0gPT0gbmlsIHRoZW4KZGF0YVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKdGFibGUuaW5zZXJ0KGRhdGFUYWJsZS5yb3dzW2pdLCBkZXZpY2VUYWJsZS5yb3dzW2pdWzFdKQplbmQKZW5kCnJldHVybiBnZW5lcmF0ZVRhYmxlKGRhdGFUYWJsZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZC4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQpsb2NhbCByZXNwLCBlcnIgPSB0cy5nZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmlmIGVyciB+PSBuaWwgdGhlbgplcnJvcihlcnIpCmVuZAppZiByZXNwID09IG5pbCB0aGVuCmVycm9yKCJubyByZWNvcmRzIGZvdW5kIikKZW5kCnJldHVybiByZXNwCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZW5lcmF0ZSBOb3J0aFN0YXIgdGFibGUuCmZ1bmN0aW9uIGdlbmVyYXRlVGFibGUodGFibGUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnRhYmxlKHRhYmxlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIHNwbGl0IGRlbGltaXRlZCBzdHJpbmdzIGludG8gbHVhIHRhYmxlLgpmdW5jdGlvbiBzcGxpdChzdHIsIGRlbGltKQppZiBzdHJpbmcuZmluZChzdHIsIGRlbGltKSA9PSBuaWwgdGhlbgpyZXR1cm4geyBzdHIgfQplbmQKbG9jYWwgYXJyYXkgPSB7fQpsb2NhbCBpbmRleCA9IDEKZm9yIHZhbHVlIGluIHN0cmluZy5nbWF0Y2goc3RyLCAiKFteIiAuLiBkZWxpbSAuLiAiXSspIikgZG8KYXJyYXlbaW5kZXhdID0gdmFsdWUKaW5kZXggPSBpbmRleCArIDEKZW5kCnJldHVybiBhcnJheQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gY2hlY2sgaWYgc3RyaW5nIGlzIGVtcHR5LgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"discreteBarChart"}}}}',
        2,
        'Kg7KjMOxx8EEzYJCrMoCVzfMCDo='
    );

// 2. MultiBar Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        545ace6a-5c8d-40ae-a559-7541d2234aca,
        'Multiple Devices - MultiBar Graph',
        'Template used to generate a MultiBar Graph for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"MultiBar Graph","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIHN1YmZpZWxkID0gY29udGV4dC5BcmdzWyJzdWJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJkZXZpY2VJZHMgaXMgZW1wdHkuIikKZW5kCmlmIGlzRW1wdHkoZmllbGQpIHRoZW4KZXJyb3IoImZpZWxkIGlzIGVtcHR5LiIpCmVuZAotLSBQYXJzZSBkZXZpY2UgaWRzLgpsb2NhbCBpZHMgPSBzcGxpdChkZXZpY2VJZHMsICIsIikKbG9jYWwgbmFtZXMgPSBzcGxpdChkZXZpY2VOYW1lcywgIiwiKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZSBmb3IgZGV2aWNlKHMpCmZvciBpID0gMSwgI2lkcywgMSBkbwpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhpZHNbaV0sIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIGNvbHVtbk5hbWUgPSBkZXZpY2VUYWJsZS5jb2x1bW5zWzFdCmlmIG5hbWVzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBuYW1lc1tpXQplbmQKZGF0YVRhYmxlLmNvbHVtbnNbaV0gPSBjb2x1bW5OYW1lCmZvciBqID0gMSwgI2RldmljZVRhYmxlLnJvd3MsIDEgZG8KaWYgZGF0YVRhYmxlLnJvd3Nbal0gPT0gbmlsIHRoZW4KZGF0YVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKdGFibGUuaW5zZXJ0KGRhdGFUYWJsZS5yb3dzW2pdLCBkZXZpY2VUYWJsZS5yb3dzW2pdWzFdKQplbmQKZW5kCnJldHVybiBnZW5lcmF0ZVRhYmxlKGRhdGFUYWJsZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZC4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQpsb2NhbCByZXNwLCBlcnIgPSB0cy5nZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmlmIGVyciB+PSBuaWwgdGhlbgplcnJvcihlcnIpCmVuZAppZiByZXNwID09IG5pbCB0aGVuCmVycm9yKCJubyByZWNvcmRzIGZvdW5kIikKZW5kCnJldHVybiByZXNwCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZW5lcmF0ZSBOb3J0aFN0YXIgdGFibGUuCmZ1bmN0aW9uIGdlbmVyYXRlVGFibGUodGFibGUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnRhYmxlKHRhYmxlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIHNwbGl0IGRlbGltaXRlZCBzdHJpbmdzIGludG8gbHVhIHRhYmxlLgpmdW5jdGlvbiBzcGxpdChzdHIsIGRlbGltKQppZiBzdHJpbmcuZmluZChzdHIsIGRlbGltKSA9PSBuaWwgdGhlbgpyZXR1cm4geyBzdHIgfQplbmQKbG9jYWwgYXJyYXkgPSB7fQpsb2NhbCBpbmRleCA9IDEKZm9yIHZhbHVlIGluIHN0cmluZy5nbWF0Y2goc3RyLCAiKFteIiAuLiBkZWxpbSAuLiAiXSspIikgZG8KYXJyYXlbaW5kZXhdID0gdmFsdWUKaW5kZXggPSBpbmRleCArIDEKZW5kCnJldHVybiBhcnJheQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gY2hlY2sgaWYgc3RyaW5nIGlzIGVtcHR5LgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"multiBarChart"}}}}',
        2,
        'Kg7KjMOxx8EEzYJCrMoCVzfMCDo='

    );

// 3. Line Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        58190552-3b6d-4779-8a78-fcabab3e827f,
        'Multiple Devices - Line Graph',
        'Template used to generate a Line Graph for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"Line Graph","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIHN1YmZpZWxkID0gY29udGV4dC5BcmdzWyJzdWJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJkZXZpY2VJZHMgaXMgZW1wdHkuIikKZW5kCmlmIGlzRW1wdHkoZmllbGQpIHRoZW4KZXJyb3IoImZpZWxkIGlzIGVtcHR5LiIpCmVuZAotLSBQYXJzZSBkZXZpY2UgaWRzLgpsb2NhbCBpZHMgPSBzcGxpdChkZXZpY2VJZHMsICIsIikKbG9jYWwgbmFtZXMgPSBzcGxpdChkZXZpY2VOYW1lcywgIiwiKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZSBmb3IgZGV2aWNlKHMpCmZvciBpID0gMSwgI2lkcywgMSBkbwpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhpZHNbaV0sIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIGNvbHVtbk5hbWUgPSBkZXZpY2VUYWJsZS5jb2x1bW5zWzFdCmlmIG5hbWVzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBuYW1lc1tpXQplbmQKZGF0YVRhYmxlLmNvbHVtbnNbaV0gPSBjb2x1bW5OYW1lCmZvciBqID0gMSwgI2RldmljZVRhYmxlLnJvd3MsIDEgZG8KaWYgZGF0YVRhYmxlLnJvd3Nbal0gPT0gbmlsIHRoZW4KZGF0YVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKdGFibGUuaW5zZXJ0KGRhdGFUYWJsZS5yb3dzW2pdLCBkZXZpY2VUYWJsZS5yb3dzW2pdWzFdKQplbmQKZW5kCnJldHVybiBnZW5lcmF0ZVRhYmxlKGRhdGFUYWJsZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZC4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQpsb2NhbCByZXNwLCBlcnIgPSB0cy5nZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmlmIGVyciB+PSBuaWwgdGhlbgplcnJvcihlcnIpCmVuZAppZiByZXNwID09IG5pbCB0aGVuCmVycm9yKCJubyByZWNvcmRzIGZvdW5kIikKZW5kCnJldHVybiByZXNwCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZW5lcmF0ZSBOb3J0aFN0YXIgdGFibGUuCmZ1bmN0aW9uIGdlbmVyYXRlVGFibGUodGFibGUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnRhYmxlKHRhYmxlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIHNwbGl0IGRlbGltaXRlZCBzdHJpbmdzIGludG8gbHVhIHRhYmxlLgpmdW5jdGlvbiBzcGxpdChzdHIsIGRlbGltKQppZiBzdHJpbmcuZmluZChzdHIsIGRlbGltKSA9PSBuaWwgdGhlbgpyZXR1cm4geyBzdHIgfQplbmQKbG9jYWwgYXJyYXkgPSB7fQpsb2NhbCBpbmRleCA9IDEKZm9yIHZhbHVlIGluIHN0cmluZy5nbWF0Y2goc3RyLCAiKFteIiAuLiBkZWxpbSAuLiAiXSspIikgZG8KYXJyYXlbaW5kZXhdID0gdmFsdWUKaW5kZXggPSBpbmRleCArIDEKZW5kCnJldHVybiBhcnJheQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gY2hlY2sgaWYgc3RyaW5nIGlzIGVtcHR5LgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"lineChart"}}}}',
        2,
        'Kg7KjMOxx8EEzYJCrMoCVzfMCDo='

    );

// 4. Pie Chart Graph
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        2338daad-fe8f-4fbc-9cd7-2dbb5a085449,
        'Multiple Devices - Pie Chart Graph',
        'Template used to generate a Pie Chart Graph for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"Pie Chart Graph","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIHN1YmZpZWxkID0gY29udGV4dC5BcmdzWyJzdWJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJkZXZpY2VJZHMgaXMgZW1wdHkuIikKZW5kCmlmIGlzRW1wdHkoZmllbGQpIHRoZW4KZXJyb3IoImZpZWxkIGlzIGVtcHR5LiIpCmVuZAotLSBQYXJzZSBkZXZpY2UgaWRzLgpsb2NhbCBpZHMgPSBzcGxpdChkZXZpY2VJZHMsICIsIikKbG9jYWwgbmFtZXMgPSBzcGxpdChkZXZpY2VOYW1lcywgIiwiKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZSBmb3IgZGV2aWNlKHMpCmZvciBpID0gMSwgI2lkcywgMSBkbwpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhpZHNbaV0sIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIGNvbHVtbk5hbWUgPSBkZXZpY2VUYWJsZS5jb2x1bW5zWzFdCmlmIG5hbWVzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBuYW1lc1tpXQplbmQKZGF0YVRhYmxlLmNvbHVtbnNbaV0gPSBjb2x1bW5OYW1lCmZvciBqID0gMSwgI2RldmljZVRhYmxlLnJvd3MsIDEgZG8KaWYgZGF0YVRhYmxlLnJvd3Nbal0gPT0gbmlsIHRoZW4KZGF0YVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKdGFibGUuaW5zZXJ0KGRhdGFUYWJsZS5yb3dzW2pdLCBkZXZpY2VUYWJsZS5yb3dzW2pdWzFdKQplbmQKZW5kCnJldHVybiBnZW5lcmF0ZVRhYmxlKGRhdGFUYWJsZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZC4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQpsb2NhbCByZXNwLCBlcnIgPSB0cy5nZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmlmIGVyciB+PSBuaWwgdGhlbgplcnJvcihlcnIpCmVuZAppZiByZXNwID09IG5pbCB0aGVuCmVycm9yKCJubyByZWNvcmRzIGZvdW5kIikKZW5kCnJldHVybiByZXNwCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZW5lcmF0ZSBOb3J0aFN0YXIgdGFibGUuCmZ1bmN0aW9uIGdlbmVyYXRlVGFibGUodGFibGUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnRhYmxlKHRhYmxlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIHNwbGl0IGRlbGltaXRlZCBzdHJpbmdzIGludG8gbHVhIHRhYmxlLgpmdW5jdGlvbiBzcGxpdChzdHIsIGRlbGltKQppZiBzdHJpbmcuZmluZChzdHIsIGRlbGltKSA9PSBuaWwgdGhlbgpyZXR1cm4geyBzdHIgfQplbmQKbG9jYWwgYXJyYXkgPSB7fQpsb2NhbCBpbmRleCA9IDEKZm9yIHZhbHVlIGluIHN0cmluZy5nbWF0Y2goc3RyLCAiKFteIiAuLiBkZWxpbSAuLiAiXSspIikgZG8KYXJyYXlbaW5kZXhdID0gdmFsdWUKaW5kZXggPSBpbmRleCArIDEKZW5kCnJldHVybiBhcnJheQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gY2hlY2sgaWYgc3RyaW5nIGlzIGVtcHR5LgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"pieChart"}}}}',
        2,
        'Kg7KjMOxx8EEzYJCrMoCVzfMCDo='

    );

// 5. Data Table
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        64ffc83a-efa1-431f-a08c-a8602ba96648,
        'Multiple Devices - Data Table',
        'Template used to generate a Data Table for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"Data Table","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":"","subfield":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIHN1YmZpZWxkID0gY29udGV4dC5BcmdzWyJzdWJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJkZXZpY2VJZHMgaXMgZW1wdHkuIikKZW5kCmlmIGlzRW1wdHkoZmllbGQpIHRoZW4KZXJyb3IoImZpZWxkIGlzIGVtcHR5LiIpCmVuZAotLSBQYXJzZSBkZXZpY2UgaWRzLgpsb2NhbCBpZHMgPSBzcGxpdChkZXZpY2VJZHMsICIsIikKbG9jYWwgbmFtZXMgPSBzcGxpdChkZXZpY2VOYW1lcywgIiwiKQpsb2NhbCBkYXRhVGFibGUgPSB7Y29sdW1ucyA9IHt9LCByb3dzID0ge319Ci0tIEdldCBmaWVsZCB2YWx1ZSBmb3IgZGV2aWNlKHMpCmZvciBpID0gMSwgI2lkcywgMSBkbwpsb2NhbCBkZXZpY2VUYWJsZSA9IGdldERldmljZUZpZWxkcyhpZHNbaV0sIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIGNvbHVtbk5hbWUgPSBkZXZpY2VUYWJsZS5jb2x1bW5zWzFdCmlmIG5hbWVzW2ldIH49IG5pbCB0aGVuCmNvbHVtbk5hbWUgPSBuYW1lc1tpXQplbmQKZGF0YVRhYmxlLmNvbHVtbnNbaV0gPSBjb2x1bW5OYW1lCmZvciBqID0gMSwgI2RldmljZVRhYmxlLnJvd3MsIDEgZG8KaWYgZGF0YVRhYmxlLnJvd3Nbal0gPT0gbmlsIHRoZW4KZGF0YVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKdGFibGUuaW5zZXJ0KGRhdGFUYWJsZS5yb3dzW2pdLCBkZXZpY2VUYWJsZS5yb3dzW2pdWzFdKQplbmQKZW5kCnJldHVybiBnZW5lcmF0ZVRhYmxlKGRhdGFUYWJsZSkKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdldCBUaGluZ1NwYWNlIGRldmljZSBmaWVsZC4KZnVuY3Rpb24gZ2V0RGV2aWNlRmllbGRzKGlkLCBmaWVsZCwgc3ViZmllbGQsIGxpbWl0KQpsb2NhbCByZXNwLCBlcnIgPSB0cy5nZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmlmIGVyciB+PSBuaWwgdGhlbgplcnJvcihlcnIpCmVuZAppZiByZXNwID09IG5pbCB0aGVuCmVycm9yKCJubyByZWNvcmRzIGZvdW5kIikKZW5kCnJldHVybiByZXNwCmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byBnZW5lcmF0ZSBOb3J0aFN0YXIgdGFibGUuCmZ1bmN0aW9uIGdlbmVyYXRlVGFibGUodGFibGUpCmxvY2FsIG91dCwgZXJyID0gb3V0cHV0LnRhYmxlKHRhYmxlKQppZihlcnIgfj0gbmlsKSB0aGVuCmVycm9yKGVycikKZW5kCnJldHVybiBvdXQKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIHNwbGl0IGRlbGltaXRlZCBzdHJpbmdzIGludG8gbHVhIHRhYmxlLgpmdW5jdGlvbiBzcGxpdChzdHIsIGRlbGltKQppZiBzdHJpbmcuZmluZChzdHIsIGRlbGltKSA9PSBuaWwgdGhlbgpyZXR1cm4geyBzdHIgfQplbmQKbG9jYWwgYXJyYXkgPSB7fQpsb2NhbCBpbmRleCA9IDEKZm9yIHZhbHVlIGluIHN0cmluZy5nbWF0Y2goc3RyLCAiKFteIiAuLiBkZWxpbSAuLiAiXSspIikgZG8KYXJyYXlbaW5kZXhdID0gdmFsdWUKaW5kZXggPSBpbmRleCArIDEKZW5kCnJldHVybiBhcnJheQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gY2hlY2sgaWYgc3RyaW5nIGlzIGVtcHR5LgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true,"visualization":{"type":"","parameters":{"type":"table"}}}}',
        2,
        'Kg7KjMOxx8EEzYJCrMoCVzfMCDo='

    );

// 6. Map
INSERT INTO account.templates
    (
        id,
        name,
        description,
        type,
        data,
        published,
        hash
    ) values
    (
        89178262-c6cc-41dd-ba5e-34a1414abc8d,
        'Map',
        'Template used to generate a Map for a field on one or more ThingSpace devices.',
        'cell',
        '{"name":"Map","input":{"type":"Code","language":"lua","arguments":{"deviceIds":"","deviceNames":"","field":"","limit":""},"entryPoint":"main","body":"LS0KLS0gQ29weXJpZ2h0IDIwMTYgVmVyaXpvbiBMYWJvcmF0b3JpZXMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi0tIFNlZSBwcm92aWRlZCBMSUNFTlNFIGZpbGUgZm9yIHVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlLgotLQpsb2NhbCB0cyA9IHJlcXVpcmUoInRzIikKbG9jYWwgb3V0cHV0ID0gcmVxdWlyZSgibnNPdXRwdXQiKQpmdW5jdGlvbiBtYWluKCkKLS0gR2V0IGlucHV0IGFyZ3VtZW50cwpsb2NhbCBkZXZpY2VJZHMgPSBjb250ZXh0LkFyZ3NbImRldmljZUlkcyJdCmxvY2FsIGRldmljZU5hbWVzID0gY29udGV4dC5BcmdzWyJkZXZpY2VOYW1lcyJdCmxvY2FsIGZpZWxkID0gY29udGV4dC5BcmdzWyJmaWVsZCJdCmxvY2FsIGxpbWl0ID0gdG9udW1iZXIoY29udGV4dC5BcmdzWyJsaW1pdCJdKQotLSBWZXJpZnkgcmVxdWlyZWQgaW5wdXRzCmlmIGlzRW1wdHkoZGV2aWNlSWRzKSB0aGVuCmVycm9yKCJFbXB0eSBkZXZpY2UgaWRzLiIpCmVuZAppZiBpc0VtcHR5KGZpZWxkKSB0aGVuCmVycm9yKCJJbnZhbGlkIGRldmljZSBmaWVsZC4iKQplbmQKLS0gUGFyc2UgZGV2aWNlIGlkcy4KbG9jYWwgaWRzID0gc3BsaXQoZGV2aWNlSWRzLCAiLCIpCmxvY2FsIG5hbWVzID0gc3BsaXQoZGV2aWNlTmFtZXMsICIsIikKbG9jYWwgbWFwID0gewp0eXBlID0gIk1hcCIsCmNlbnRlciA9IHtsYXRpdHVkZSA9IDEsIGxvbmdpdHVkZSA9IDF9LAp6b29tID0gMTAsCml0ZW1zID0ge30KfQpsb2NhbCBjZW50ZXJTZXQgPSBmYWxzZQotLSBHZXQgZmllbGQgdmFsdWUgZm9yIGRldmljZShzKQpmb3IgaSA9IDEsICNpZHMsIDEgZG8KbG9jYWwgZGV2aWNlVGFibGUgPSBnZXREZXZpY2VGaWVsZHMoaWRzW2ldLCBmaWVsZCwgIiIsIGxpbWl0KQpsb2NhbCBjb2x1bW5OYW1lID0gIkRldmljZS0iIC4uIGlkc1tpXQppZiBuYW1lc1tpXSB+PSBuaWwgdGhlbgpjb2x1bW5OYW1lID0gbmFtZXNbaV0KZW5kCmxvY2FsIGl0ZW0gPSB7CmxhYmVsID0gY29sdW1uTmFtZSwKbG9jYXRpb25zID0ge30KfQpsb2NhbCBrZXlzID0ge30KZm9yIGkgPSAxLCAjZGV2aWNlVGFibGUuY29sdW1ucywgMSBkbwprZXlzW2RldmljZVRhYmxlLmNvbHVtbnNbaV1dID0gaQplbmQKaWYga2V5c1sibGF0aXR1ZGUiXSA9PSBuaWwgb3Iga2V5c1sibG9uZ2l0dWRlIl0gPT0gbmlsIHRoZW4KZXJyb3IoIk9uIG9yIG1vcmUgZGV2aWNlIGZpZWxkIGlzIG5vdCBhIGxvY2F0aW9uIG9iamVjdC4iKQplbmQKZm9yIGogPSAxLCAjZGV2aWNlVGFibGUucm93cywgMSBkbwppZiBkZXZpY2VUYWJsZS5yb3dzW2pdID09IG5pbCB0aGVuCmRldmljZVRhYmxlLnJvd3Nbal0gPSB7fQplbmQKbG9jYWwgbG9jYXRpb24gPSB7CmxhdGl0dWRlICA9IGRldmljZVRhYmxlLnJvd3Nbal1ba2V5c1sibGF0aXR1ZGUiXV0sCmxvbmdpdHVkZSA9IGRldmljZVRhYmxlLnJvd3Nbal1ba2V5c1sibG9uZ2l0dWRlIl1dCn0KLS0gRm9yIG5vdywgdXNlIHRoZSBmaXJzdCBsb2NhdGlvbiBhcyBjZW50ZXIuCmlmIGNlbnRlclNldCA9PSBmYWxzZSB0aGVuCm1hcC5jZW50ZXIubGF0aXR1ZGUgPSBsb2NhdGlvbi5sYXRpdHVkZQptYXAuY2VudGVyLmxvbmdpdHVkZSA9IGxvY2F0aW9uLmxvbmdpdHVkZQplbmQKdGFibGUuaW5zZXJ0KGl0ZW0ubG9jYXRpb25zLCBsb2NhdGlvbikKZW5kCnRhYmxlLmluc2VydChtYXAuaXRlbXMsIGl0ZW0pCmVuZApyZXR1cm4gZ2VuZXJhdGVNYXAobWFwKQplbmQKLS0gSGVscGVyIGZ1bmN0aW9uIHVzZWQgdG8gZ2V0IFRoaW5nU3BhY2UgZGV2aWNlIGZpZWxkLgpmdW5jdGlvbiBnZXREZXZpY2VGaWVsZHMoaWQsIGZpZWxkLCBzdWJmaWVsZCwgbGltaXQpCmxvY2FsIHJlc3AsIGVyciA9IHRzLmdldERldmljZUZpZWxkcyhpZCwgZmllbGQsIHN1YmZpZWxkLCBsaW1pdCkKaWYgZXJyIH49IG5pbCB0aGVuCmVycm9yKGVycikKZW5kCmlmIHJlc3AgPT0gbmlsIHRoZW4KZXJyb3IoIm5vIHJlY29yZHMgZm91bmQiKQplbmQKcmV0dXJuIHJlc3AKZW5kCi0tIEhlbHBlciBmdW5jdGlvbiB1c2VkIHRvIGdlbmVyYXRlIE5vcnRoU3RhciBtYXAuCmZ1bmN0aW9uIGdlbmVyYXRlTWFwKG1hcCkKbG9jYWwgb3V0LCBlcnIgPSBvdXRwdXQubWFwKG1hcCkKaWYoZXJyIH49IG5pbCkgdGhlbgplcnJvcihlcnIpCmVuZApyZXR1cm4gb3V0CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdXNlZCB0byB0cmFuc2xhdGUgY29tbWEgc2VwYXJhdGVkIHN0cmluZ3MgaW50byBsdWEgdGFibGUuCmZ1bmN0aW9uIHNwbGl0KHN0ciwgZGVsaW0pCmlmIHN0cmluZy5maW5kKHN0ciwgZGVsaW0pID09IG5pbCB0aGVuCnJldHVybiB7IHN0ciB9CmVuZApsb2NhbCBhcnJheSA9IHt9CmxvY2FsIGluZGV4ID0gMQpmb3IgdmFsdWUgaW4gc3RyaW5nLmdtYXRjaChzdHIsICIoW14iIC4uIGRlbGltIC4uICJdKykiKSBkbwphcnJheVtpbmRleF0gPSB2YWx1ZQppbmRleCA9IGluZGV4ICsgMQplbmQKcmV0dXJuIGFycmF5CmVuZAotLSBIZWxwZXIgZnVuY3Rpb24gdG8gY2hlY2sgZW1wdHkgc3RyaW5nLgpmdW5jdGlvbiBpc0VtcHR5KHN0cikKcmV0dXJuIHN0ciA9PSBuaWwgb3Igc3RyID09ICIiCmVuZA==","timeout":180},"output":{"lastExecution":"0001-01-01T00:00:00Z"},"settings":{"status":"","showCode":false,"showOutput":true}}',
        2,
        'J6Z1HGuu1JYapv+DSwKCV70cFvI='
    );
