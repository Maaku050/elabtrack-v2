import js from '@eslint/js'
import tseslint from 'typescript-eslint'
import globals from 'globals'
import hooks from 'eslint-plugin-react-hooks'
export default [
 {ignores:['dist/**','node_modules/**']},
 js.configs.recommended,
 ...tseslint.configs.recommended,
 {files:['src/**/*.{ts,tsx}'],languageOptions:{globals:globals.browser},plugins:{'react-hooks':hooks},rules:{'react-hooks/rules-of-hooks':'error','react-hooks/exhaustive-deps':'warn'}},
 {files:['scripts/**/*.mjs','*.{js,ts,mjs}'],languageOptions:{globals:globals.node}},
]
