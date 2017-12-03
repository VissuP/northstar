export class Category {
    constructor(public id: number,
                public name: string,
                public description: string) {
    }
}

export class CategoryOption {
    constructor(public id: number,
                public categoryId: number,
                public name: string,
                public description: string,
                public value: string) {
    }
}

export class CategoryOptionField {
    constructor(public id: number,
                public categoryOptionId: number,
                public name: string,
                public description: string,
                public value: string) {
    }
}