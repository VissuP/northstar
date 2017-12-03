/**
 * Copyright 2016 Verizon Laboratories. All rights reserved.
 * See provided LICENSE file for use of this source code.
 */



// Load the implementations that should be tested
import {Transformation} from "./transformation.model";

describe("Test encode() method", () => {
    let pojo = {
        code: {
            type: "type",
            value: "U291cmNlIGNvZGUgYm9keSBoZXJlIQ==",
        },
        description: "Test transformation",
        name: "test",
    };

    let transformation = new Transformation(pojo);
    let encodedTransformation = transformation.encode();

    it("Validate constructor", () => {
        expect(transformation.code.value).toContain("Source code body here!");
    });
}
)