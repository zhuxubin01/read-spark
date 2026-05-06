import Combine
import Foundation

final class HomeViewModel: ObservableObject {
    @Published var articles: [ArticleSummary] = []
    @Published var isLoading = false
    @Published var error: String?

    private var cancellables = Set<AnyCancellable>()
    private let repository = ArticleRepository.shared

    init() {
        loadDailyArticles()
    }

    func loadDailyArticles() {
        isLoading = true
        error = nil

        repository.fetchDailyArticles()
            .sink(receiveCompletion: { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.error = String(describing: error)
                }
            }, receiveValue: { [weak self] articles in
                self?.articles = articles
            })
            .store(in: &cancellables)
    }
}
