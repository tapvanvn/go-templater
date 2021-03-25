A templater need:
- Manage the relationship with other template base on url or path.
- Manage cache to optimize speed. ex: the rendered result, the loaded template...

And should:
- Support language natural editor. Template define language should using tags if woking with html file, js code if working with javascript, so the editor will display correctly make it easy to use.
- Support language natural using context.
- Has an interface to interactive with hosted language.


Implement 
- Templater is a singleton object that serves all demain about template.
- The refId of the template. Each teamplate has a refID that help to find the real file. syntax: namespace:relative_path_to_namespace. "/" is use for directory separate.
- The language type of template file is base on it's extension. ex: if the template file name is .html or .htm then it's content will be treat as html.
- A language can be relative with other language base on reality. ex: a html file can contain javascript or css content.
- In javascript, templater is named _$ in template.

Supported:
- HTML(+ss)