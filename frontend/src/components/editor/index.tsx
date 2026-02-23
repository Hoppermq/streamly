import { useCallback, useState } from "react"
import CodeMirror from '@uiw/react-codemirror'
import { sql } from "@codemirror/lang-sql"

export const CodeEditor = () => {
  const [content, setContent] = useState<string>('')
  const onChange = useCallback((v: string, _) => {
    setContent(v)
  }, [])

  return <CodeMirror
    value={content}
    height="50px"
    extensions={[sql()]}
    onChange={onChange}
  />
}
