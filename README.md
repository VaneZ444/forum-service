# forum-service
 посмотреть про перетягивание ошибок в хендлер из юзкейса, посмотреть про протягивание запроса внутрь юзкейса
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
rpc ListComments(ListCommentsRequest) returns (ListCommentsResponse); может падать пагинация
rpc UpdateComment(UpdateCommentRequest) returns (CommentResponse); он существует но проверить нечем
rpc DeleteComment(DeleteCommentRequest) returns (Empty); он существует но проверить нечем

// Tags
rpc CreateTag(CreateTagRequest) returns (TagResponse);
rpc GetTag(GetTagRequest) returns (TagResponse);
rpc ListTags(ListTagsRequest) returns (ListTagsResponse);

// Tag-Post Relationships
rpc AddTagToPost(AddTagToPostRequest) returns (Empty);
rpc RemoveTagFromPost(RemoveTagFromPostRequest) returns (Empty);
rpc ListTagsByPost(ListTagsByPostRequest) returns (ListTagsResponse);
rpc ListPostsByTag(ListPostsByTagRequest) returns (ListPostsResponse);

// Search
rpc Search(SearchRequest) returns (SearchResponse);
