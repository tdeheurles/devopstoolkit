import re


class TemplateFile:
    def __init__(self, template_file: str, template_content: str):
        self.path = template_file
        self.content = template_content
        matches = re.findall("<%.*?%>", template_content, re.MULTILINE)
        self.required_params = [
            match.translate({ord(c): None for c in '<% >'})
            for match in matches
        ]
