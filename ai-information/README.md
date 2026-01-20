# AI 工具学习指南

> 本项目详细介绍 Chatbox、Streamlit、Ollama、LangChain 四个 AI 相关工具的原理、用法和它们之间的关系。

## 目录结构

```
ai-information/
├── README.md           # 本文件 - 总览和快速入门
├── chatbox.md          # Chatbox 详解
├── streamlit.md        # Streamlit 详解
├── ollama.md           # Ollama 详解
├── langchain.md        # LangChain 详解
└── comparison.md       # 四者区别与联系（重点推荐）
```

## 快速了解

### 一句话介绍

| 工具 | 定位 | 类比 |
|------|------|------|
| **Ollama** | 本地 AI 模型运行时 | 发动机 - 提供动力 |
| **LangChain** | AI 应用开发框架 | 变速箱 - 传递和调节动力 |
| **Streamlit** | Web 界面开发框架 | 车身 - 让人能坐进去 |
| **Chatbox** | 桌面 AI 客户端 | 整车 - 买来就能开 |

### 它们的关系

```
┌─────────────────────────────────────────────────────┐
│                    用户界面层                        │
│        Chatbox (现成的)    Streamlit (自己做的)     │
└────────────────────────┬────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│                    应用框架层                        │
│                    LangChain                        │
│            (RAG、Agent、Memory 等能力)              │
└────────────────────────┬────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│                     模型层                           │
│       Ollama (本地)        OpenAI/Claude (云端)     │
└─────────────────────────────────────────────────────┘
```

## 学习路线

### 按技能水平

| 水平 | 推荐组合 | 说明 |
|------|----------|------|
| **零基础** | Chatbox + Ollama | 下载安装即可使用，无需编程 |
| **会点 Python** | Streamlit + Ollama | 20 行代码做 AI 应用 |
| **熟悉 Python** | Streamlit + LangChain + Ollama | 构建复杂 AI 应用 |

### 按需求场景

| 需求 | 推荐方案 |
|------|----------|
| 只想聊天 | Chatbox + Ollama |
| 做简单 Demo | Streamlit + Ollama |
| 知识库问答 | Streamlit + LangChain + Ollama |
| 让 AI 用工具 | LangChain Agent + Ollama |
| 数据完全私密 | 任意组合，用 Ollama 本地模型 |

## 快速开始

### 1. 安装 Ollama（基础）

```bash
# macOS / Linux
curl -fsSL https://ollama.com/install.sh | sh

# 运行模型
ollama run llama3.2
```

### 2. 安装 Chatbox（零编程使用）

从 [chatboxai.app](https://chatboxai.app) 下载，配置 Ollama 地址即可使用。

### 3. 安装 Streamlit + LangChain（开发使用）

```bash
pip install streamlit langchain langchain-ollama langchain-community
```

### 4. 运行示例应用

```python
# app.py
import streamlit as st
from langchain_ollama import ChatOllama

st.title("AI 助手")
llm = ChatOllama(model="llama3.2")

if prompt := st.chat_input("输入问题"):
    st.chat_message("user").write(prompt)
    response = llm.invoke(prompt)
    st.chat_message("assistant").write(response.content)
```

```bash
streamlit run app.py
```

## 文档详情

### [chatbox.md](./chatbox.md)
- Chatbox 是什么
- 工作原理和架构
- 与同类产品对比
- 使用场景和快速上手

### [streamlit.md](./streamlit.md)
- Streamlit 核心概念
- 执行模型（脚本重跑机制）
- 常用组件速览
- 与 Gradio、Flask 对比

### [ollama.md](./ollama.md)
- 本地模型运行原理
- GGUF 格式和量化
- 支持的模型列表
- API 使用和集成

### [langchain.md](./langchain.md)
- 为什么需要 LangChain
- 核心组件详解（Chain、RAG、Agent、Memory）
- 代码示例
- 与 LlamaIndex 对比

### [comparison.md](./comparison.md) ⭐ 推荐
- 四者的区别与联系
- 协作方式和数据流
- 选择指南
- 完整代码示例

## 核心概念图

```
┌──────────────────────────────────────────────────────────────┐
│                                                               │
│  问题: "AI 怎么跑？"          →  Ollama (运行模型)           │
│  问题: "AI 怎么更强？"        →  LangChain (扩展能力)        │
│  问题: "怎么做界面？"         →  Streamlit (开发界面)        │
│  问题: "不想编程怎么用？"     →  Chatbox (现成界面)          │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

## 推荐学习顺序

```
1. Ollama       → 理解 AI 模型如何在本地运行
       ↓
2. Chatbox      → 体验完整的 AI 对话流程
       ↓
3. Streamlit    → 学会用 Python 做 Web 界面
       ↓
4. LangChain    → 掌握 RAG、Agent 等高级能力
       ↓
5. 整合项目     → 用全套技术栈构建完整应用
```

## 实战案例：从本地到云端的架构迁移

> 使用阿里云平台提供的云上模型构建聊天机器人，技术架构的主要改动：

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           架构迁移总结                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  改动 1: 模型管理层                                                      │
│  ──────────────────                                                     │
│  原方案：Ollama (本地模型)                                               │
│  新方案：LangChain 框架                                                  │
│  说明：使用 LangChain 代替 Ollama 完成对模型的管理和使用，               │
│       基于 LangChain 提供 Python 调用支持                               │
│                                                                          │
│  改动 2: 算力来源                                                        │
│  ──────────────────                                                     │
│  原方案：本地 CPU/GPU                                                    │
│  新方案：阿里云百炼平台 (通义千问系列模型)                               │
│  说明：使用云端模型提供算力支持，无需本地硬件                            │
│                                                                          │
│  改动 3: 前端界面                                                        │
│  ──────────────────                                                     │
│  方案：保持不变，继续使用 Streamlit 原有代码                             │
│  说明：前端与后端解耦，切换模型来源不影响界面                            │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 迁移前后对比

| 层级 | 迁移前 | 迁移后 |
|------|--------|--------|
| **前端界面** | Streamlit | Streamlit（不变） |
| **应用框架** | 直接调用 Ollama | LangChain |
| **模型来源** | Ollama 本地模型 | 阿里云百炼（通义千问） |
| **算力** | 本地 CPU/GPU | 云端算力 |

### 迁移后的代码示例

```python
# 使用阿里云通义千问的示例
import streamlit as st
from langchain_community.llms import Tongyi

st.title("AI 助手 (通义千问)")

# 使用阿里云百炼平台的通义千问模型
llm = Tongyi(model="qwen-turbo")  # 需要设置 DASHSCOPE_API_KEY

if prompt := st.chat_input("输入问题"):
    st.chat_message("user").write(prompt)
    response = llm.invoke(prompt)
    st.chat_message("assistant").write(response)
```

### 核心理解

```
这个案例说明了 LangChain 的核心价值：

┌──────────────────────────────────────────────────────────────┐
│                                                               │
│   Streamlit (界面)  ──►  LangChain (框架)  ──►  模型 (可替换) │
│                               │                               │
│                               ├──► Ollama (本地)             │
│                               ├──► 通义千问 (阿里云)          │
│                               ├──► OpenAI (国外)             │
│                               └──► 其他模型...               │
│                                                               │
│   关键：LangChain 提供统一接口，切换模型只需改一行代码        │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

## 相关资源

- [Ollama 官网](https://ollama.com)
- [Chatbox 官网](https://chatboxai.app)
- [Streamlit 文档](https://docs.streamlit.io)
- [LangChain 文档](https://python.langchain.com)
- [阿里云百炼平台](https://bailian.console.aliyun.com)

---

**提示**：建议先阅读 [comparison.md](./comparison.md)，快速理解四者的关系，再根据需要深入学习各个工具。
