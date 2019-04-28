'use strict';

const fs = require("fs");
const path = require("path");
const Swagger = require('swagger-client');

function readSchemaFromLocalFile() {
    return JSON.parse(fs.readFileSync('./schema/petstore.json', 'utf8'));
}

function readJsonFromFile(path) {
    return JSON.parse(fs.readFileSync(path));
}

const readImageFile = (filename) => {
    return fs.readFileSync(filename);
    // convert binary data to base64 encoded string
    // return new Buffer(bitmap).toString('base64');
};

describe('petstore rest api integration tests', () => {
    let apiClient = {};

    test('api client can be generated from local schema file', () => {
        expect.assertions(2);

        return Swagger({
            url: 'test',
            spec: readSchemaFromLocalFile(),
            requestInterceptor: req => {
                // console.log(req);
                return req;
            }
        }).then((client) => {
            expect(client).not.toBeNull();
            expect(client.apis.default).not.toBeNull();
            apiClient = client;
        });
    });

    test('create user', () => {
        var input = readJsonFromFile("./__tests__/user_create.json");
        expect.assertions(2);

        return apiClient.apis.user.createUser({body: input}).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('create users with array', () => {
        var input = readJsonFromFile("./__tests__/user_createWithArray.json");
        expect.assertions(2);

        return apiClient.apis.user.createUsersWithArrayInput({body: input}).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('create users with list', () => {
        var input = readJsonFromFile("./__tests__/user_createWithList.json");
        expect.assertions(2);

        return apiClient.apis.user.createUsersWithListInput({body: input}).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('user login', () => {
        expect.assertions(4);

        return apiClient.apis.user.loginUser({
            username: "username1",
            password: "string",
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
            expect(resp.headers["x-rate-limit"]).toBe("100");
            expect(resp.headers["x-expires-after"]).toBe("3600");
        }).catch((err) => {
            console.log(err);
        });
    });

    test('user logout', () => {
        expect.assertions(2);

        return apiClient.apis.user.logoutUser({}).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('user get by name', () => {
        expect.assertions(3);

        return apiClient.apis.user.getUserByName({
            username: "username1",
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
            expect(resp.body).toEqual(readJsonFromFile("./__tests__/username1.json"))
        }).catch((err) => {
            console.log(err);
        });
    });

    test('user update by username', () => {
        expect.assertions(2);

        return apiClient.apis.user.updateUser({
            username: "username1",
            body: readJsonFromFile("./__tests__/user_update.json"),
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('user delete by username', () => {
        expect.assertions(2);

        return apiClient.apis.user.deleteUser({
            username: "username1",
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(204);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet add', () => {
        expect.assertions(2);

        return apiClient.apis.pet.addPet({
            body: readJsonFromFile("./__tests__/pet_add.json")
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet add 2', () => {
        expect.assertions(2);

        return apiClient.apis.pet.addPet({
            body: readJsonFromFile("./__tests__/pet_add2.json")
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet add 3', () => {
        expect.assertions(2);

        return apiClient.apis.pet.addPet({
            body: readJsonFromFile("./__tests__/pet_add3.json")
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet update', () => {
        expect.assertions(2);

        return apiClient.apis.pet.updatePet({
            body: readJsonFromFile("./__tests__/pet_update.json")
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet find by status', () => {
        expect.assertions(3);

        return apiClient.apis.pet.findPetsByStatus({
            status: "available,pending"
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
            expect(resp.body.length).toBe(2);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet find by id', () => {
        expect.assertions(3);

        return apiClient.apis.pet.getPetById({
            petId: 1,
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
            expect(resp.body).toEqual(readJsonFromFile("./__tests__/pet_add2.json"));
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet update with form', () => {
        expect.assertions(2);

        return apiClient.apis.pet.updatePetWithForm({
            petId: 1,
            name: "pet1",
            status: "sold",
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('pet delete by id', () => {
        expect.assertions(2);

        return apiClient.apis.pet.deletePet({
            petId: 1,
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(204);
        }).catch((err) => {
            console.log(err);
        });
    });

    // test('pet upload image by id', () => {
    //     expect.assertions(2);
    //
    //     const filepath = path.normalize(__dirname + "/test.jpg");
    //     const image = fs.createReadStream("__tests__/test.jpg");
    //     const data = {
    //         'file': image
    //     };
    //     // const formData = new FormData();
    //     // formData.append('file', image, "test.jpg");
    //     return apiClient.apis.pet.uploadFile({
    //         petId: 2,
    //         file: fs.createReadStream("__tests__/test.jpg"),
    //         additionalMetadata: "test",
    //     }, {
    //         requestInterceptor: req => {
    //             req.headers["accept"] = "application/json";
    //             req.headers["Content-Type"] = "multipart/form-data";
    //             return req;
    //         }
    //     }).then(resp => {
    //         
    //         expect(resp).not.toBeNull();
    //         expect(resp.status).toBe(200);
    //     }).catch((err) => {
    //         console.log(err);
    //     });
    // });

    test('store get inventory', () => {
        expect.assertions(2);

        return apiClient.apis.store.getInventory({}).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('store place order', () => {
        expect.assertions(2);

        return apiClient.apis.store.placeOrder({
            body: readJsonFromFile("./__tests__/store_order.json")
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('store get order', () => {
        expect.assertions(2);

        return apiClient.apis.store.getOrderById({
            orderId: 0,
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(200);
        }).catch((err) => {
            console.log(err);
        });
    });

    test('store delete order', () => {
        expect.assertions(2);

        return apiClient.apis.store.deleteOrder({
            orderId: 0,
        }).then(resp => {

            expect(resp).not.toBeNull();
            expect(resp.status).toBe(204);
        }).catch((err) => {
            console.log(err);
        });
    });
});
