import Editor, { type OnMount } from "@monaco-editor/react";

export function EditorPane({
  code,
  onChange,
}: {
  code: string;
  onChange: (v: string | undefined) => void;
}) {
  const handleMount: OnMount = (editor) => {
    editor.focus();
  };

  return (
    <div className="editor-pane">
      <Editor
        height="100%"
        defaultLanguage="go"
        theme="vs-dark"
        value={code}
        onChange={onChange}
        onMount={handleMount}
        options={{
          minimap: { enabled: false },
          fontSize: 13,
          tabSize: 4,
          wordWrap: "on",
          automaticLayout: true,
        }}
        loading={<em>starting go language server…</em>}
      />
    </div>
  );
}
