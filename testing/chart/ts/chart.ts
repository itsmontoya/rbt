import * as parser from "node_modules/utils/parser";

const xmlns = "http://www.w3.org/2000/svg"

export class Chart {
	private ta: HTMLTextAreaElement;
	private chart: HTMLElement;

	constructor() {
		this.ta = document.createElement("textarea");
		this.ta.onchange = () => this.handleTextAreaChange();

		this.chart = document.createElement("chart");

		document.body.appendChild(this.ta);
		document.body.appendChild(this.chart);
	}

	private handleTextAreaChange() {
		const obj = JSON.parse(this.ta.value);
		if (obj === null) {
			console.error("invalid JSON");
			return;
		}

		this.Process(obj);
	}

	Process(obj: {}) {
		this.chart.innerHTML = "";
		const db = new DebugBlock(this.chart, obj);
		console.log("Processed!", db);
	}
}

export function NewDebugBlock(parent: HTMLElement | SVGElement, obj: {} | null): Child {
	if (obj === null) {
		return null;
	}

	return new DebugBlock(parent, obj)
}

export class DebugBlock {
	private e: HTMLElement;
	private cir: HTMLElement;
	private cs: HTMLElement;

	readonly childType: number;
	readonly color: number;
	readonly parent: string;
	readonly key: string;
	readonly children: Child[];

	constructor(parent: HTMLElement | SVGElement, obj: {} | null) {
		this.e = document.createElement("block");
		this.cir = document.createElement("circle");
		this.cs = document.createElement("children");

		this.e.appendChild(this.cir);
		this.e.appendChild(this.cs);
		parent.appendChild(this.e);

		if(!!obj){
			const p = new parser.Parser(obj);
			this.childType = p.Number("childType");
			this.color = p.Number("color");
			this.parent = p.String("parent");
			this.key = p.String("key");

			const children = p.Array("children");
			this.children = new Array(2);
			this.children[0] = new DebugBlock(this.cs, children[0]);
			this.children[1] = new DebugBlock(this.cs, children[1]);
		} else {
			this.key = "null";
			this.color = 0;
		}

		this.cir.textContent = this.key;
		console.log("Uh", this.color);
		switch(this.color){
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

export type Child = DebugBlock | null;

function newSVG(type: string): SVGElement {
	const e = document.createElementNS(xmlns, type);
	return e;
}

function newCircle(cx: number, cy: number, r: number): SVGElement {
	const e = newSVG("circle");
	e.setAttribute("cx", cx.toString());
	e.setAttribute("cy", cy.toString());
	e.setAttribute("r", r.toString());
	return e;
}