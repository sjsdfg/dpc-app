// See config-examples/template.js for full list of options and
// their descriptions.
module.exports = {
    name: "DPC API",
    description: "Data at the Point of Care Bulk FHIR API",
    authType: "backend-services", // backend-services|client-credentials|none
    jwksAuth: false,
    baseURL: `${process.env.BASE_URL}/api/v1`,
    tokenEndpoint: `${process.env.BASE_URL}/v1/Token/auth`,
    systemExportEndpoint: "",
    grouExportEndpoint: `/Group/${process.env.GROUP_ID}/$export`,
    patientExportEndpoint: `/Patient/${process.env.PATIENT_ID}/$everything`,
    clientId: process.env.CLIENT_ID,
    strictSSL: true,
    requiresAuth: true,
    fastestResource: "Group",
    groupId: process.env.GROUP_ID,
    sinceParam: "_since",
    scope: "system/*.*",
    privateKey: process.env.PRIVATE_KEY
};
