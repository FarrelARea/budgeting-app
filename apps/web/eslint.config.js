/* eslint-disable */
const base = require("@budgeting/eslint-config/base");
module.exports = [...base, { ignores: [".next/**", "out/**", "node_modules/**", "next-env.d.ts"] }];
