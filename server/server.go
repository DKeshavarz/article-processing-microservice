package server

import (
	"context"
	"errors"
	"io"
	"log"

	"article-processing-microservice/proto"
	"article-processing-microservice/tagextractor"

)

type ArticleServer struct {
	proto.UnimplementedArticleServiceServer
}

func NewArticleServer() *ArticleServer {
	return &ArticleServer{}
}

func (s *ArticleServer) ProcessSingleArticle(ctx context.Context, req *proto.ProcessArticleRequest) (*proto.ProcessArticleResponse, error) {
	log.Printf("Processing single article: %s", req.Article.Title)

	tags, err := tagextractor.ExtractTags(req.Article.Body, int(req.N))
	if err != nil {
		log.Printf("Error extracting tags: %v", err)
		return nil, err
	}

	return &proto.ProcessArticleResponse{Tags: tags}, nil
}

func (s *ArticleServer) ProcessArticles(stream proto.ArticleService_ProcessArticlesServer) error {
	log.Println("Starting concurrent bidirectional streaming processing...")

	for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }

        if err != nil {
            return errors.New("error recive stream")
        }

  
        go func(request *proto.ProcessArticleRequest) {
            tags, err := tagextractor.ExtractTags(request.Article.Body, int(request.N))
            if err != nil {
                log.Printf("Error extracting tags: %v", err)
            }

            if err := stream.Send(&proto.ProcessArticleResponse{Tags: tags}); err != nil {
                log.Printf("Error sending response: %v", err)
            }
        }(req)
    }
}

//TODO:
func (s *ArticleServer) GetTopTags(ctx context.Context, req *proto.GetTopTagsRequest) (*proto.GetTopTagsResponse, error) {

	return &proto.GetTopTagsResponse{Tags: []string{"nigger"}}, nil
}
