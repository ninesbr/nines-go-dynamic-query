# 🧠 dynamicquery

**dynamicquery** é uma biblioteca leve, extensível e plugável para construir **consultas dinâmicas em APIs RESTful** baseadas em GORM (Go ORM).  
Ela permite aplicar filtros (`filter`), ordenações (`sort`), seleções (`select`) e paginação diretamente via **query string**, com suporte a aliases e fallback automático CamelCase → snake_case.

---

## 🚀 Instalação

```bash
go get github.com/ninesbr/dynamicquery@latest