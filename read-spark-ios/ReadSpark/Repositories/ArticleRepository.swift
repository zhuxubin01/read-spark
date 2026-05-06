import Combine
import CoreData
import Foundation

final class ArticleRepository {
    static let shared = ArticleRepository()

    private let api = APIService.shared
    private let context = CoreDataStack.shared.container.viewContext

    func fetchDailyArticles() -> AnyPublisher<[ArticleSummary], APIError> {
        api.getDailyArticles()
            .map { $0.articles }
            .eraseToAnyPublisher()
    }

    func fetchArticle(id: UUID) -> AnyPublisher<Article, APIError> {
        api.getArticle(id: id)
    }

    func fetchArticles(category: String? = nil,
                       difficulty: String? = nil,
                       page: Int = 1) -> AnyPublisher<ArticleListResponse, APIError> {
        api.getArticles(category: category, difficulty: difficulty, page: page)
    }

    func cache(_ article: ArticleSummary) {
        let entity = ArticleEntity(context: context)
        entity.id = article.id
        entity.title = article.title
        entity.summary = article.summary
        entity.category = article.category
        entity.difficulty = article.difficulty
        entity.wordCount = Int32(article.wordCount)
        entity.isPremium = article.isPremium
        entity.updatedAt = Date()
        CoreDataStack.shared.save()
    }
}
