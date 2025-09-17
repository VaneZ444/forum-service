# forum-service
TODO: посмотреть про перетягивание ошибок в хендлер из юзкейса, посмотреть про протягивание запроса внутрь юзкейса
rpc CreateCategory(CreateCategoryRequest) returns (CategoryResponse); Работает
rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse); работет
rpc GetCategory(GetCategoryRequest) returns (CategoryResponse); работет
rpc UpdateCategory(UpdateCategoryRequest) returns (CategoryResponse); работет
rpc DeleteCategory(DeleteCategoryRequest) returns (Empty); Работает

// Topics
rpc CreateTopic(CreateTopicRequest) returns (TopicResponse); Работает
rpc GetTopic(GetTopicRequest) returns (TopicResponse); Работает
rpc ListTopics(ListTopicsRequest) returns (ListTopicsResponse); Работает
rpc UpdateTopic(UpdateTopicRequest) returns (TopicResponse); Работает
rpc DeleteTopic(DeleteTopicRequest) returns (Empty); Работает

// Posts
rpc CreatePost(CreatePostRequest) returns (PostResponse); Работает
rpc GetPost(GetPostRequest) returns (PostResponse); Работает
rpc ListPosts(ListPostsRequest) returns (ListPostsResponse); Работает
rpc UpdatePost(UpdatePostRequest) returns (PostResponse); Работает
rpc DeletePost(DeletePostRequest) returns (Empty); Работает

// Comments
rpc CreateComment(CreateCommentRequest) returns (CommentResponse); Работает
rpc GetComment(GetCommentRequest) returns (CommentResponse); Работает
rpc ListComments(ListCommentsRequest) returns (ListCommentsResponse); Работает
rpc UpdateComment(UpdateCommentRequest) returns (CommentResponse); Работает
rpc DeleteComment(DeleteCommentRequest) returns (Empty); Работает

// Tags
rpc CreateTag(CreateTagRequest) returns (TagResponse); Работает
rpc GetTag(GetTagRequest) returns (TagResponse); Работает
rpc ListTags(ListTagsRequest) returns (ListTagsResponse); Работает

// Tag-Post Relationships
rpc AddTagToPost(AddTagToPostRequest) returns (Empty); Работает
rpc RemoveTagFromPost(RemoveTagFromPostRequest) returns (Empty); Работает
rpc ListTagsByPost(ListTagsByPostRequest) returns (ListTagsResponse); Работает
rpc ListPostsByTag(ListPostsByTagRequest) returns (ListPostsResponse); Работает

// Search
rpc Search(SearchRequest) returns (SearchResponse); Работает
