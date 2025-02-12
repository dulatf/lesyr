export interface User {
    id: string;
    username: string;
    email: string;
    avatarUrl: string;
  }
  
  export interface Feed {
    id: string;
    userId: string;
    url: string;
    title: string;
    description: string;
  }
  
  export interface Article {
    id: string;
    feedId: string;
    title: string;
    content: string;
    url: string;
    publishedAt: string;
    read: boolean;
    createdAt: string;
  }