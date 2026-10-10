package store

import (
	"context"
	"fmt"

	"portfolio/backend/internal/schema"
)

// demoContent is the placeholder content the site ships with. It is inserted
// once, only when every content table is empty, so a fresh install looks
// like the finished design instead of a blank page.
var demoContent = map[string][]map[string]string{
	"education": {
		{"degree": "B.Sc. Computer Science", "school": "University Name", "years": "2020 – 2024", "note": "GPA 3.9/4.0 · Thesis: Efficient retrieval for long-context LLMs"},
	},
	"skills": {
		{"name": "AI", "items": "PyTorch\nTransformers\nvLLM\nRAG\nLangChain\nCUDA basics"},
		{"name": "Backend", "items": "Python\nGo\nFastAPI\nPostgreSQL\nRedis\ngRPC\nKafka"},
		{"name": "Tools", "items": "Docker\nKubernetes\nAWS\nGitHub Actions\nLinux"},
	},
	"profiles": {
		{"group_name": "Professional", "label": "GitHub", "href": "https://github.com/yourhandle", "icon": "github"},
		{"group_name": "Professional", "label": "LinkedIn", "href": "https://linkedin.com/in/yourhandle", "icon": "linkedin"},
		{"group_name": "Professional", "label": "X", "href": "https://x.com/yourhandle", "icon": "x"},
		{"group_name": "Academic", "label": "Google Scholar", "href": "https://scholar.google.com/citations?user=XXXX", "icon": "scholar"},
		{"group_name": "Academic", "label": "ORCID", "href": "https://orcid.org/0000-0000-0000-0000", "icon": "orcid"},
		{"group_name": "Academic", "label": "ResearchGate", "href": "https://researchgate.net/profile/yourhandle", "icon": "research"},
		{"group_name": "Coding", "label": "LeetCode", "href": "https://leetcode.com/u/yourhandle", "icon": "code"},
		{"group_name": "Coding", "label": "Codeforces", "href": "https://codeforces.com/profile/yourhandle", "icon": "trophy"},
		{"group_name": "Coding", "label": "CodeChef", "href": "https://codechef.com/users/yourhandle", "icon": "chef"},
		{"group_name": "Coding", "label": "HackerRank", "href": "https://hackerrank.com/profile/yourhandle", "icon": "terminal"},
		{"group_name": "Coding", "label": "Kaggle", "href": "https://kaggle.com/yourhandle", "icon": "chart"},
	},
	"experience": {
		{"title": "Software Engineer", "company": "Company A · ML Platform", "period": "Jun 2024 – Present"},
		{"title": "Research Assistant", "company": "Systems Lab, University Name", "period": "Sep 2023 – May 2024"},
	},
	"projects": {
		{"slug": "fast-rag", "title": "FastRAG", "summary": "Retrieval pipeline that cut p95 query latency 40% with hybrid BM25 + dense search.", "tags": "AI,Backend", "stack": "Python\nFastAPI\npgvector\nRedis", "featured": "1", "github_url": "https://github.com/yourhandle/fast-rag", "live_url": "https://fast-rag.example.com"},
		{"slug": "kv-cache-scheduler", "title": "KV-Cache Scheduler", "summary": "Request scheduler for LLM serving that raised throughput 1.8× at equal latency.", "tags": "AI,Research", "stack": "Python\nPyTorch\nvLLM\nCUDA", "featured": "1", "github_url": "https://github.com/yourhandle/kv-scheduler", "live_url": "https://kv-scheduler.example.com"},
		{"slug": "event-ledger", "title": "Event Ledger", "summary": "Append-only event store in Go handling 50k writes/s with exactly-once consumers.", "tags": "Backend", "stack": "Go\nPostgreSQL\nKafka\ngRPC", "featured": "1", "github_url": "https://github.com/yourhandle/event-ledger", "live_url": "https://event-ledger.example.com"},
		{"slug": "ml-feature-store", "title": "Feature Store", "summary": "Online/offline feature store that removed training–serving skew for three models.", "tags": "AI,Backend", "stack": "Python\nSpark\nRedis\nAirflow", "featured": "1", "github_url": "https://github.com/yourhandle/feature-store", "live_url": "https://feature-store.example.com"},
		{"slug": "cli-profiler", "title": "pyprof", "summary": "Sampling profiler CLI for Python async code with flame-graph output.", "tags": "Backend", "stack": "Python\nRust", "featured": "0", "github_url": "https://github.com/yourhandle/pyprof", "live_url": "https://pyprof.example.com"},
	},
	"papers": {
		{"title": "Length-Aware Scheduling for Memory-Efficient LLM Inference", "authors": "Rahat Almas\nCo Author\nAdvisor Name", "venue": "Workshop on ML Systems, NeurIPS", "year": "2025", "status": "published", "pdf_url": "https://arxiv.org/abs/0000.00000", "code_url": "https://github.com/yourhandle/kv-scheduler"},
		{"title": "Hybrid Retrieval for Long-Context Question Answering", "authors": "Rahat Almas\nAdvisor Name", "venue": "arXiv preprint", "year": "2025", "status": "preprint", "pdf_url": "https://arxiv.org/abs/0000.00001"},
		{"title": "Efficient Retrieval for Long-Context LLMs", "authors": "Rahat Almas", "venue": "B.Sc. Thesis, University Name", "year": "2024", "status": "published", "pdf_url": "/thesis.pdf"},
	},
}

// SeedDemoContent inserts demoContent if, and only if, every resource table
// is empty. It returns true when it seeded.
func (s *Store) SeedDemoContent(ctx context.Context) (bool, error) {
	for i := range schema.Resources {
		n, err := s.CountRecords(ctx, &schema.Resources[i])
		if err != nil {
			return false, err
		}
		if n > 0 {
			return false, nil
		}
	}
	for i := range schema.Resources {
		res := &schema.Resources[i]
		for _, row := range demoContent[res.Key] {
			values := make(map[string]string, len(res.Fields))
			for _, f := range res.Fields {
				v, problem := f.Clean(row[f.Key])
				if problem != "" {
					return false, fmt.Errorf("demo %s: %s", res.Key, problem)
				}
				values[f.Key] = v
			}
			if _, err := s.CreateRecord(ctx, res, values); err != nil {
				return false, fmt.Errorf("seed %s: %w", res.Key, err)
			}
		}
	}
	return true, nil
}
