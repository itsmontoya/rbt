define("node_modules/utils/parser", ["require", "exports"], function (require, exports) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    class Parser {
        constructor(obj) {
            this.obj = obj;
        }
        String(key) {
            const v = this.obj[key];
            if (typeof v !== "string") {
                if (v !== undefined) {
                    console.error("invalid value, expected string and received:", v);
                }
                return "";
            }
            return v;
        }
        Number(key) {
            const v = this.obj[key];
            if (typeof v !== "number") {
                if (v !== undefined) {
                    console.error("invalid value, expected number and received:", v);
                }
                return 0;
            }
            return v;
        }
        Boolean(key) {
            const v = this.obj[key];
            if (typeof v !== "boolean") {
                if (v !== undefined) {
                    console.error("invalid value, expected boolean and received:", v);
                }
                return false;
            }
            return v;
        }
        Object(key) {
            const v = this.obj[key];
            if (typeof v !== "object") {
                if (v !== undefined) {
                    console.error("invalid value, expected object and received:", v);
                }
                return {};
            }
            return v;
        }
        Array(key) {
            const v = this.obj[key];
            if (!Array.isArray(v)) {
                if (v !== undefined) {
                    console.error("invalid value, expected array and received:", v);
                }
                return [];
            }
            return v;
        }
    }
    exports.Parser = Parser;
});
define("chart", ["require", "exports", "node_modules/utils/parser"], function (require, exports, parser) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    const xmlns = "http://www.w3.org/2000/svg";
    class Chart {
        constructor() {
            this.ta = document.createElement("textarea");
            this.ta.onchange = () => this.handleTextAreaChange();
            this.chart = document.createElement("chart");
            document.body.appendChild(this.ta);
            document.body.appendChild(this.chart);
        }
        handleTextAreaChange() {
            const obj = JSON.parse(this.ta.value);
            if (obj === null) {
                console.error("invalid JSON");
                return;
            }
            this.Process(obj);
        }
        Process(obj) {
            this.chart.innerHTML = "";
            const db = new DebugBlock(this.chart, obj);
            console.log("Processed!", db);
        }
    }
    exports.Chart = Chart;
    function NewDebugBlock(parent, obj) {
        if (obj === null) {
            return null;
        }
        return new DebugBlock(parent, obj);
    }
    exports.NewDebugBlock = NewDebugBlock;
    class DebugBlock {
        constructor(parent, obj) {
            this.e = document.createElement("block");
            this.cir = document.createElement("circle");
            this.cs = document.createElement("children");
            this.e.appendChild(this.cir);
            this.e.appendChild(this.cs);
            parent.appendChild(this.e);
            if (!!obj) {
                const p = new parser.Parser(obj);
                this.childType = p.Number("childType");
                this.color = p.Number("color");
                this.parent = p.String("parent");
                this.key = p.String("key");
                const children = p.Array("children");
                this.children = new Array(2);
                this.children[0] = new DebugBlock(this.cs, children[0]);
                this.children[1] = new DebugBlock(this.cs, children[1]);
            }
            else {
                this.key = "null";
                this.color = 0;
            }
            this.cir.textContent = this.key;
            console.log("Uh", this.color);
            switch (this.color) {
                case 0:
                    this.cir.classList.add("black");
                    break;
                case 1:
                    this.cir.classList.add("red");
                    break;
                case 2:
                    this.cir.classList.add("doubleBlack");
                    break;
            }
        }
    }
    exports.DebugBlock = DebugBlock;
    function newSVG(type) {
        const e = document.createElementNS(xmlns, type);
        return e;
    }
    function newCircle(cx, cy, r) {
        const e = newSVG("circle");
        e.setAttribute("cx", cx.toString());
        e.setAttribute("cy", cy.toString());
        e.setAttribute("r", r.toString());
        return e;
    }
});
