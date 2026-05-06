import Foundation

struct Article: Codable, Identifiable {
    let id: UUID
    let title: String
    let summary: String?
    let content: String
    let translation: String?
    let category: String
    let difficulty: String
    let wordCount: Int
    let audioUrl: String?
    let coverImage: String?
    let isPremium: Bool
    let publishedAt: String?
}

struct ArticleSummary: Codable, Identifiable {
    let id: UUID
    let title: String
    let summary: String?
    let category: String
    let difficulty: String
    let wordCount: Int
    let coverImage: String?
    let isPremium: Bool
    let publishedAt: String?
}

struct TokenPair: Codable {
    let accessToken: String
    let refreshToken: String
    let expiresIn: Int
}

struct ReadingProgress: Codable, Identifiable {
    let id: UUID
    let articleId: UUID
    let position: Int
    let percentage: Double
    let lastReadAt: String
}
