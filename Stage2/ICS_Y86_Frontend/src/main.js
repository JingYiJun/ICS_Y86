let current_p = 0;  // 记录json文件运行到哪了
    let json;
    let code;
    let max_len;
    let started = false;
    let table;
    let timer;
    let res;

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("src/backend.wasm"), go.importObject).then((result) => go.run(result.instance));

    let inputFiles = document.querySelectorAll('input[type=file]');
    for (let i = 0; i < inputFiles.length; i++) {
        let inputFile = inputFiles[i];
        inputFile.onchange = function () {
            let reader = new FileReader();
            let file = this.files[0];
            reader.readAsText(file);
            reader.onload = function () {
                code = this.result;
                document.getElementById("code").innerHTML = this.result;
                document.getElementById("subTitle").innerHTML = file.name;
                res = goRun(code);

                if (res.error !== "") {
                    alert("yo文件解析失败，请检查您的文件格式");
                } else {
                    json = JSON.parse(res.result);
                    // error = res.error;
                    updateTable(json, 0);
                    let but = document.getElementById("pause");
                    but.innerText = "暂停";
                    Reset();
                }
            }
        }
    }

    function vanish() {
        document.getElementById("disappear").remove();
        document.getElementById("disappear1").remove();
    }

    function updateTable(json, p) {
        table = document.getElementById("myTable");
        let len = table.rows.length
        if (len !== 3) {
            for (let i = 0; i < (len - 3); i = i + 1) {
                table.deleteRow(table.rows.length - 1);
            }
        }
        let first = true;
        let mem_num = 0;
        for (const property in json[p].MEM) {
            mem_num = mem_num + 1;
        }
        for (const property in json[p].MEM) {
            let newRow = table.insertRow(table.rows.length);
            if (first) {
                newRow.insertCell(0).innerHTML = "MEM";
                let firstCell = table.rows[table.rows.length - 1].cells[0];
                firstCell.setAttribute("rowspan", mem_num.toString());
                newRow.insertCell(1).innerHTML = property;
                newRow.insertCell(2).innerHTML = 0;
                let cell = table.rows[table.rows.length - 1].cells[2];
                let id = property.toString();
                cell.setAttribute("id", id);
                first = false;
                continue;
            }
            newRow.insertCell(0).innerHTML = property;
            newRow.insertCell(1).innerHTML = 0;
            let cell = table.rows[table.rows.length - 1].cells[1];
            let id = property.toString();
            cell.setAttribute("id", id);
        }
    }


    function update(p) {
        document.getElementById("v_rax").innerHTML = (json[p].REG.rax >>> 0).toString(16);
        document.getElementById("v_rbx").innerHTML = (json[p].REG.rbx >>> 0).toString(16);
        document.getElementById("v_rbp").innerHTML = (json[p].REG.rbp >>> 0).toString(16);
        document.getElementById("v_rsp").innerHTML = (json[p].REG.rsp >>> 0).toString(16);
        document.getElementById("v_rdi").innerHTML = (json[p].REG.rdi >>> 0).toString(16);
        document.getElementById("v_rsi").innerHTML = (json[p].REG.rsi >>> 0).toString(16);
        document.getElementById("v_r8").innerHTML = (json[p].REG.r8 >>> 0).toString(16);
        document.getElementById("v_r9").innerHTML = (json[p].REG.r9 >>> 0).toString(16);
        document.getElementById("v_r10").innerHTML = (json[p].REG.r10 >>> 0).toString(16);
        document.getElementById("v_r11").innerHTML = (json[p].REG.r11 >>> 0).toString(16);
        document.getElementById("v_r12").innerHTML = (json[p].REG.r12 >>> 0).toString(16);
        document.getElementById("v_r13").innerHTML = (json[p].REG.r13 >>> 0).toString(16);
        document.getElementById("v_r14").innerHTML = (json[p].REG.r14 >>> 0).toString(16);
        document.getElementById("v_rcx").innerHTML = (json[p].REG.rcx >>> 0).toString(16);
        document.getElementById("v_rdx").innerHTML = (json[p].REG.rdx >>> 0).toString(16);

        document.getElementById("v_pc").innerHTML = "0x0" + (json[p].PC).toString(16);
        document.getElementById("v_zf").innerHTML = (json[p].CC.ZF).toString(16);
        document.getElementById("v_sf").innerHTML = (json[p].CC.SF).toString(16);
        document.getElementById("v_of").innerHTML = (json[p].CC.OF).toString(16);

        let stat = "";
        if (json[p].STAT === 1) stat = "AOK";
        else if (json[p].STAT === 2) stat = "HLT";
        else if (json[p].STAT === 3) stat = "ADR";
        else if (json[p].STAT === 4) stat = "INS";
        document.getElementById("v_stat").innerHTML = stat;

        for (const property in json[p].MEM) {
            let id = property.toString();
            document.getElementById(id).innerHTML = ((json[p].MEM[id]));
        }
    }


    function TravelJSON() {
        if (current_p >= max_len) {
            clearTimeout(timer);
        }
        started = true;
        updateTable(json, current_p);
        update(current_p);
        highlightLine(json[current_p].PC);
        current_p++;
        timer = setTimeout("TravelJSON()", 2000);
    }

    // 单步执行
    function SingleStep() {
        if (started === false) {
            current_p = 0;
            started = true;
        } else current_p = current_p + 1;

        if (current_p >= max_len) return;
        updateTable(json, current_p);
        update(current_p);
        highlightLine(json[current_p].PC);
    }

    function Pause() {
        clearTimeout(timer);
    }

    function Continue() {
        TravelJSON();
    }

    let btn = document.getElementById("pause");
    btn.addEventListener('click', function () {
        let text = btn.innerText;
        if (text === "暂停") {
            Pause();
            btn.innerText = "继续";
        } else {
            Continue();
            btn.innerText = "暂停";
        }
    });

    // 重置，回到yo文件和json文件都回到开头
    function Reset() {
        max_len = Object.keys(json).length;
        current_p = 0;
        updateTable(json, 0);
        highlightLine(0);
        started = false;
        // update(current_p);
        document.getElementById("v_rax").innerHTML = 0;
        document.getElementById("v_rbx").innerHTML = 0;
        document.getElementById("v_rbp").innerHTML = 0;
        document.getElementById("v_rsp").innerHTML = 0;
        document.getElementById("v_rdi").innerHTML = 0;
        document.getElementById("v_rsi").innerHTML = 0;
        document.getElementById("v_r8").innerHTML = 0;
        document.getElementById("v_r9").innerHTML = 0;
        document.getElementById("v_r10").innerHTML = 0;
        document.getElementById("v_r11").innerHTML = 0;
        document.getElementById("v_r12").innerHTML = 0;
        document.getElementById("v_r13").innerHTML = 0;
        document.getElementById("v_r14").innerHTML = 0;
        document.getElementById("v_rcx").innerHTML = 0;
        document.getElementById("v_rdx").innerHTML = 0;

        document.getElementById("v_pc").innerHTML = 0;
        document.getElementById("v_zf").innerHTML = 0;
        document.getElementById("v_sf").innerHTML = 0;
        document.getElementById("v_of").innerHTML = 0;
        document.getElementById("v_stat").innerHTML = "";

        for (const property in json[0].MEM) {
            let id = property.toString();
            document.getElementById(id).innerHTML = 0;
        }
    }

    function highlightLine(pc) {
        let lines = code.split("\n");
        // 利用Javascript的map()方法遍历代码行数组
        let highlightLines = lines.map(function (code, i) {
            let line = lines[i].split(":");
            // 获取行号，将十六进制字符串转为十进制整数
            let num = parseInt(line[0], 16);
            // 如果当前行是需要高亮的行，则将其添加<span>标签
            if (num === pc) {
                return "<span class='highlight'>" + lines[i] + "</span>"
            }
            // 否则直接返回
            return lines[i];
        })
        highlightLines = highlightLines.join("\n");
        document.getElementById("code").innerHTML = highlightLines;
    }