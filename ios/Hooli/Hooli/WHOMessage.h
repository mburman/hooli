//
//  WHOMessage.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface WHOMessage : NSObject

- (instancetype) initWithMessage:(NSString* )message Author:(NSString* )author Distance:(NSString* )distance Location:(NSString* )location;
@property (nonatomic, strong) NSString* message;
@property (nonatomic, strong) NSString* author;
@property (nonatomic, strong) NSString* distance;
@property (nonatomic, strong) NSString* location;

@end
