import XCTest
@testable import ReadSpark

final class ModelDecodingTests: XCTestCase {
    func testDecodeArticle() throws {
        let json = """
        {
          "id": "11111111-2222-3333-4444-555555555555",
          "title": "A1 English News",
          "summary": "Short summary",
          "content": "Full article content",
          "translation": "完整译文",
          "category": "news",
          "difficulty": "A1",
          "wordCount": 120,
          "audioUrl": "https://example.com/audio.mp3",
          "coverImage": "https://example.com/cover.jpg",
          "isPremium": false,
          "publishedAt": "2026-05-06T12:00:00Z"
        }
        """.data(using: .utf8)!

        let article = try JSONDecoder().decode(Article.self, from: json)

        XCTAssertEqual(article.title, "A1 English News")
        XCTAssertEqual(article.category, "news")
        XCTAssertEqual(article.wordCount, 120)
        XCTAssertFalse(article.isPremium)
    }

    func testDecodeArticleListResponse() throws {
        let json = """
        {
          "articles": [
            {
              "id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
              "title": "Item 1",
              "summary": "S1",
              "category": "tech",
              "difficulty": "B1",
              "wordCount": 300,
              "coverImage": null,
              "isPremium": true,
              "publishedAt": "2026-05-06T12:00:00Z"
            }
          ],
          "total": 1,
          "page": 1,
          "page_size": 20
        }
        """.data(using: .utf8)!

        let response = try JSONDecoder().decode(ArticleListResponse.self, from: json)

        XCTAssertEqual(response.total, 1)
        XCTAssertEqual(response.page, 1)
        XCTAssertEqual(response.page_size, 20)
        XCTAssertEqual(response.articles.count, 1)
        XCTAssertEqual(response.articles.first?.title, "Item 1")
        XCTAssertTrue(response.articles.first?.isPremium == true)
    }

    func testDecodeTokenPairAndReadingProgress() throws {
        let tokenJSON = """
        {
          "accessToken": "access-token",
          "refreshToken": "refresh-token",
          "expiresIn": 900
        }
        """.data(using: .utf8)!

        let progressJSON = """
        {
          "id": "123e4567-e89b-12d3-a456-426614174000",
          "articleId": "123e4567-e89b-12d3-a456-426614174001",
          "position": 400,
          "percentage": 37.5,
          "lastReadAt": "2026-05-06T12:00:00Z"
        }
        """.data(using: .utf8)!

        let token = try JSONDecoder().decode(TokenPair.self, from: tokenJSON)
        let progress = try JSONDecoder().decode(ReadingProgress.self, from: progressJSON)

        XCTAssertEqual(token.accessToken, "access-token")
        XCTAssertEqual(token.expiresIn, 900)
        XCTAssertEqual(progress.position, 400)
        XCTAssertEqual(progress.percentage, 37.5)
    }
}
